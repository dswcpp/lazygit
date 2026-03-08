package agent

import (
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

type ReviewPhase string

const (
	PhaseReviewInit        ReviewPhase = "review_init"        // 初始化
	PhaseReviewing         ReviewPhase = "reviewing"          // 评审中
	PhaseReviewDone        ReviewPhase = "review_done"        // 完成
	PhaseReviewWaiting     ReviewPhase = "review_waiting"     // 等待用户追问
	PhaseReviewInteractive ReviewPhase = "review_interactive" // 交互式追问中
	PhaseReviewCancelled   ReviewPhase = "review_cancelled"   // 取消
	PhaseReviewError       ReviewPhase = "review_error"       // 错误
)

// CodeReviewState 代码评审状态（不可变）
type CodeReviewState struct {
	Phase      ReviewPhase
	FilePath   string
	Diff       string
	Language   string
	Focus      string // "security" | "performance" | "correctness" | ""
	Messages   []provider.Message
	Result     string
	Error      string
	StartTime  time.Time
	FinishTime time.Time

	// 交互式评审支持
	ResumeFrom     NodeID // 恢复点（用于检查点）
	UserQuestion   string // 用户追问
	ConversationID string // 会话ID（用于批量评审）
}

// WithPhase 返回新状态（不可变）
func (s CodeReviewState) WithPhase(phase ReviewPhase) CodeReviewState {
	s.Phase = phase
	return s
}

// WithResult 返回新状态（不可变）
func (s CodeReviewState) WithResult(result string) CodeReviewState {
	s.Result = result
	s.FinishTime = time.Now()
	return s
}

// WithError 返回新状态（不可变）
func (s CodeReviewState) WithError(err string) CodeReviewState {
	s.Error = err
	s.Phase = PhaseReviewError
	s.FinishTime = time.Now()
	return s
}

// WithMessages 返回新状态（深拷贝）
func (s CodeReviewState) WithMessages(messages []provider.Message) CodeReviewState {
	newMessages := make([]provider.Message, len(messages))
	copy(newMessages, messages)
	s.Messages = newMessages
	return s
}

// AppendMessage 追加消息（深拷贝）
func (s CodeReviewState) AppendMessage(msg provider.Message) CodeReviewState {
	newMessages := make([]provider.Message, len(s.Messages), len(s.Messages)+1)
	copy(newMessages, s.Messages)
	s.Messages = append(newMessages, msg)
	return s
}

// WithResumeFrom 设置恢复点
func (s CodeReviewState) WithResumeFrom(nodeID NodeID) CodeReviewState {
	s.ResumeFrom = nodeID
	return s
}

// WithUserQuestion 设置用户追问
func (s CodeReviewState) WithUserQuestion(question string) CodeReviewState {
	s.UserQuestion = question
	return s
}

// WithConversationID 设置会话ID
func (s CodeReviewState) WithConversationID(id string) CodeReviewState {
	s.ConversationID = id
	return s
}

// detectLanguage 从文件路径推断语言
func detectLanguage(filePath string) string {
	lastDot := strings.LastIndex(filePath, ".")
	if lastDot == -1 {
		return ""
	}
	ext := strings.ToLower(filePath[lastDot+1:])
	langs := map[string]string{
		"go": "Go", "ts": "TypeScript", "tsx": "TypeScript/React",
		"js": "JavaScript", "jsx": "JavaScript/React", "py": "Python",
		"rs": "Rust", "java": "Java", "c": "C", "h": "C",
		"cpp": "C++", "cc": "C++", "hpp": "C++",
		"rb": "Ruby", "php": "PHP", "swift": "Swift",
		"kt": "Kotlin", "cs": "C#", "sh": "Shell", "bash": "Shell",
		"yaml": "YAML", "yml": "YAML", "json": "JSON", "sql": "SQL",
	}
	return langs[ext]
}
