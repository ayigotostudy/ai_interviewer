# {{name}}
{{ gender }} ｜ {{ age }}岁 ｜ {{ profession }} ｜ {{ education }} ｜ {{ phone }} ｜ {{ email }}
## 自我评价
{{ self_evaluation }}

## 工作经历
{% for exp in work_experience %}
::: start
**{{ exp.company }}**

:::
**{{ exp.position }}**

:::
**{{ exp.duration }}**

::: end
{{ exp.highlights }}
{% endfor %}

## 项目经历
{% for project in projects %}
::: start
**{{ project.name }}**

:::
**{{ project.role }}**

:::
**{{ project.scope }}**

:::
**{{ project.period }}**

::: end
{{ project.achievements }}
{% endfor %}

## 专业技能
{{ technical_skills }}

## 教育背景与资质
::: start
**{{ university }}** | {{ major }} | {{ degree }}

:::
**{{ education_period }}**

::: end
{{ certifications }}