package trace

import "io"

// コード内のフローを記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
  return nil
}
