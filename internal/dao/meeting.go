package dao

import (
	"ai_jianli_go/types/model"

	"gorm.io/gorm"
)

type MeetingDAO struct {
	db *gorm.DB
}

func NewMeetingDAO(db *gorm.DB) *MeetingDAO {

	return &MeetingDAO{db: db}
}

func (dao *MeetingDAO) Create(meeting *model.Meeting) error {
	return dao.db.Create(meeting).Error
}

func (dao *MeetingDAO) Update(meeting *model.Meeting) error {
	return dao.db.Save(meeting).Error
}

func (dao *MeetingDAO) GetByID(id uint) (*model.Meeting, error) {
	var meeting model.Meeting
	err := dao.db.First(&meeting, id).Error
	return &meeting, err
}

func (dao *MeetingDAO) List() ([]model.Meeting, error) {
	var meetings []model.Meeting
	err := dao.db.Find(&meetings).Error
	return meetings, err
}

func (dao *MeetingDAO) Delete(id uint) error {
	return dao.db.Delete(&model.Meeting{}, id).Error
}

func (dao *MeetingDAO) UploadResume(candidateID uint, resume string) error {
	return dao.db.Model(&model.Meeting{}).Where("id = ?", candidateID).Update("resume", resume).Error
}

func (dao *MeetingDAO) GetResume(candidateID uint) (string, error) {
	var meeting model.Meeting
	err := dao.db.Select("resume").Where("id = ?", candidateID).First(&meeting).Error
	return meeting.Resume, err
}
