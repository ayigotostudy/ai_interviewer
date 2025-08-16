package main

import (
	"ai_jianli_go/pkg/rag"
	"context"
	"os"
	"testing"

	"github.com/cloudwego/eino/components/document"
)

func TestXxx(t *testing.T) {
	logFile, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// log := func(format string, args ...interface{}) {
	// 	fmt.Fprintf(logFile, format+"\n", args...)
	// }

	ctx := context.Background()
	rag.Init()
	// 获取文档
	docs, err := rag.GetLoader().Load(ctx, document.Source{
		URI: "./test1.pdf",
	})
	if err != nil {
		t.Error("文档解析失败：", err)
		return
	}

	// 存储
	indexer := rag.GetIndexer()
	_, err = indexer.Store(ctx, docs)

	if err != nil {
		t.Error(err)
	}

	// 	// 检索
	// 	retriever := component.GetRetriever()

	// 	// 对话模型
	// 	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
	// 		Model:   "gpt-4o",
	// 		BaseURL: "https://api.vveai.com/v1",
	// 		APIKey:  "sk-Xs5rROO2htFLGMJh407b42F505Fe4c89A8510f7608E52c2f",
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// 记忆
	// 	memory := component.NewRedisMemory(component.RedisMemoryConfig{
	// 		MaxWindowSize: 20,
	// 		RedisOptions: &redis.Options{
	// 			Addr:     "124.222.151.35:6379",
	// 			Password: "123456",
	// 		},
	// 	})

	// 	con := memory.GetConversation("1", true)
	// 	con.SetLastConversationKnowledge("java")
	// 	fmt.Println("开始面试")
	// 	for {
	// 		docs, err := retriever.Retrieve(ctx, con.GetLastConversationsKnowledge())
	// 		if err != nil {
	// 			fmt.Println("文档检索失败：", err)
	// 			return
	// 		}

	// 		// 模型问答
	// 		context := ""
	// 		// 查看文档内容
	// 		if len(docs) > 0 {
	// 			contextParts := make([]string, len(docs))
	// 			for i, doc := range docs {
	// 				contextParts[i] = fmt.Sprintf("文档片段[%d]:\n%s\n", i+1, doc.Content)
	// 			}
	// 			context = strings.Join(contextParts, "\n---\n")
	// 		}

	// 		// 创建模板，使用 FString 格式
	// 		template := prompt.FromMessages(schema.FString,
	// 			// 系统角色定义 + 核心任务说明 + 输出格式要求
	// 			schema.SystemMessage(
	// 				"你是一个专业面试官，需要完成以下任务：\n"+
	// 					"1. 基于专业知识库内容提出精准问题\n"+
	// 					"2. 对用户回答进行结构化评价（优点/不足）\n"+
	// 					"3. 针对不足点给出专业解释\n"+
	// 					"4. 根据用户回答生成1-4轮深度追问， 追问结束继续根据简历内容提问\n"+
	// 					"5. 每次回答都需要返回要问的知识点（关键词）以便后续知识库检索提问\n\n"+
	// 					"6. 追问每次只追问一道题目， 后续在根据用户回答继续追问，最多追问4轮\n"+
	// 					"当前知识库上下文：{context}\n\n"+
	// 					"当前对话记录：{history}\n\n"+
	// 					"用户简历内容:{resume}\n"+
	// 					"输出格式要求：\n"+
	// 					"- 评价使用✅和❌标识优劣点\n",
	// 			),

	// 			// 用户回答捕获
	// 			schema.UserMessage("【应聘者回答】\n{answer}"),

	// 			// 系统回答
	// 			schema.AssistantMessage(
	// 				"请按以下结构组织回答：\n"+
	// 					"1. 评价总结（含具体不足点分析）\n"+
	// 					"2. 可追问的知识点:knowledgepoint",
	// 				[]schema.ToolCall{},
	// 			),
	// 		)

	// 		var prompt map[string]any = make(map[string]any)

	// 		prompt["context"] = context
	// 		answer := ""
	// 		fmt.Scanln(&answer)
	// 		fmt.Println("🧑‍ : ", answer)
	// 		log("用户输入: %s", answer)
	// 		prompt["answer"] = answer
	// 		prompt["resume"] = `
	// 		广东石油化工学院 - 本科 - 数据科学与大数据技术
	// 2022.09 - 2026.06 - 广东
	// 综合成绩： 专业前15%；GPA:3.8/5;CET-4:480
	// <!-- - 可以删除 -->
	// 荣誉奖项： 校级一等奖学金、蓝桥杯省二
	// 语雀: https://www.yuque.com/zhongxiashier-qqvas/study
	// Github: https://github.com/zhongxia12
	// 嘉为蓝鲸-广州嘉为科技有限公司-后端研发实习生
	// 主要工作：
	// 对接腾讯ieg下的自研paas平台蓝鲸业务
	// 通过命令行参数动态化（配置解析 + 参数实时注入），解决了静态参数导致的资源利用率瓶颈，将 CPU 利用率从
	// 60% 提升至 85%，任务吞吐量提升 45%。
	// 针对多种主从数据处理，通过分步查询 + 映射缓存优化，结合策略与工厂模式解耦多模板逻辑，将数据库交互时间
	// 复杂度从 O(n) 优化为 O(1)，查询耗时降低 64%。
	// 采用 AI 驱动的多模块分析 + 黑白盒测试方法，编写测试用例，识别出大量边界条件和异常场景，将测试覆盖率从
	// 58% 提升至 83%，漏测缺陷减少 75%。
	// 通过 DDD 基础设施层重构 ，实现公共组件统一管理，建设防腐层隔离业务依赖，采用灰度迁移策略，消除代码冗
	// 余，将变更成本降低 80%，新功能接入效率提升 50%。
	// 设计项目缓存架构,并实现了MySQL-Redis-本地内存三级缓存体系，将核心接口响应时间从200ms降至30ms。
	// 🎰积分抽奖营销激励系统 2024.07 - 2024.10
	// 技术栈: 、 、 、 、 、 、 、 、
	// 项目描述: 参考掘金社区的签到抽奖和拼多多大转盘开发的抽奖系统。满足C端人群拉新、促活、留存等需求；支持签
	// 到返利、积分兑换、累计解锁等功能。为提升系统扩展和可维护性，采用多种设计模式。同时针对性优化了抽奖算法
	// 和秒杀场景，可支持单机4c16g服务器500~800TPS,抽奖接口响应时长45~100ms。
	// 技术亮点:
	// 系统架构: 采用DDD领域驱动设计，业务拆分为策略领域(奖品策略)、责任领域(参与抽奖)、奖品领域(发放奖品)、
	// 任务领域(消息发送)，并使用多种设计模式，如：组合模式、工厂模式、责任链模式和模板模式；
	// 流程设计: 使用模板模式定义抽奖、抽奖中与抽奖后各阶段流程；抽奖前进行订单创建和参数校验；抽奖中使用工
	// 厂模式组装责任链，进行黑 / 白名单和权重过滤；抽奖后利用组合模式构造决策树模型做解锁、库存和兜底奖品的
	// 进一步过滤，支持差异化营销活动。
	// 数据库路由: 使用自研组件DBRouter，基于哈希散列和斐波那契扰乱算法，依据用户 ID 对抽奖订单、中奖订单等
	// 高频率进行分库分表实现。通过Canal同步 Binlog 日志至Elasticsearch，支持高效聚合查询。
	// 接口可用性: 通过AOP编程和自定义注解结合Redission令牌桶实现动态超频次访问限流，引入Hystrix提供服务熔
	// 断能力；利用Zookeeper实现分布式动态配置降级、限流、超时熔断开关组件，提高系统灵活性。
	// 高并发处理: 使用Redis decr库存扣减，并借鉴ConcurrentHashMap的分段锁机制，设计非独占竞争的分段化兜
	// 底防超卖，通过延迟队列+定时任务异步处理保证数据最终一致性，缓解数据库并发压力。
	// 后续加一个spring ai
	// Java：熟练掌握Java基础知识，熟悉常用集合、Stream流、Lambda表达式的使用，了解Java 8新特性。
	// MySQL：熟悉掌握MySQL，理解事务、索引、日志、锁机制、MVCC等，具备一定SQL调优能力。
	// Redis：熟悉掌握Redis，理解常见数据结构、持久化机制、内存淘汰、哨兵机制，缓存故障等。
	// JUC：熟悉Java并发容器的使用及原理，对线程池、重入锁、AQS、volatile关键字有一定了解。
	// JVM：熟悉类加载机制、JVM的内存区域、垃圾回收器等；了解内存溢出、内存泄漏等问题。
	// 消息队列：熟悉RabbitMQ的使用，理解一致性、不重复消费、顺序消费、延迟消费、消息堆积等问题。
	// 框架：熟练使用SSM、Spring Boot、Spring Cloud等框架，理解IOC、AOP、Bean生命周期等。
	// 设计模式：熟悉使用常见的设计模式，如：责任链、策略、模板、工厂等模式，提高代码可复用性。
	// 🎓教育背景
	// 高水平理工科大学 卓越工程师培养计划
	// 📋实习经历
	// ✏️项目经验
	// DDD SpringBoot Redis MySQL MyBatis RabbitMQ ZooKeeper ES Canal
	// 🎯专业技能
	// 架构：了解DDD领域驱动设计，如充血模型、依赖倒置、四色建模等，了解DDD分层架构、六边形架构等。
	// 其它：熟练使用Maven、Git、Docker等开发管理工具，有Prometheus + Grafana监控系统部署经验。

	// 		`
	// 		prompt["history"] = con.String()
	// 		// 使用模板生成消息
	// 		messages, err := template.Format(ctx, prompt)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		log("", messages)
	// 		con.Append(schema.UserMessage(fmt.Sprintln(answer)))
	// 		res, err := chatModel.Generate(ctx, messages)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		fmt.Println(extractKnowledgePoint(res.Content))
	// 		con.SetLastConversationKnowledge(extractKnowledgePoint(res.Content))
	// 		con.Append(res)

	// 		fmt.Println("🤖 : ", res.Content)

	// 	}

	// }

	// func extractKnowledgePoint(input string) string {
	// 	// 查找"可追问的知识点："的位置
	// 	prefix := "可追问的知识点："
	// 	index := strings.Index(input, prefix)

	// 	if index >= 0 {
	// 		// 提取前缀之后的所有内容
	// 		result := input[index+len(prefix):]
	// 		// 去除开头可能的多余空格
	// 		return strings.TrimSpace(result)
	// 	}

	// 	// 兼容可能使用英文冒号的情况
	// 	prefix = "可追问的知识点:"
	// 	index = strings.Index(input, prefix)
	// 	if index >= 0 {
	// 		result := input[index+len(prefix):]
	// 		return strings.TrimSpace(result)
	// 	}

	//		return "" // 未找到匹配
	//	}
}
