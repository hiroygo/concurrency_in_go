# concurrency_in_go
* `Go言語による並行処理` を読んで書いてみたプログラム

## 日本語版サポートページ
* https://github.com/oreilly-japan/concurrency-in-go-support

## 原著サンプルコード
* https://github.com/kat-co/concurrency-in-go-src

## メモ
### p47
* struct{} は空構造体と呼ばれる。メモリを消費しない
* 計測などに使うとよい

### p49
* ゴルーチンがスケジュールされるタイミングにはなんの保証もない
* sync.WaitGroup の Add はできる限りゴルーチンの直前に書く
```
var wg sync.WaitGroup

// ゴルーチンの外側で Add する
// ゴルーチンの内側で Add すると Add 前に Wait が実行されて
// ゴルーチン終了の待機が行われない可能性がある
wg.Add(1)

go func() {
    defer wg.Done()
    fmt.Println("Hello")
}()

wg.Wait()
```

### p77
* チャネルの作成者の責任: チャネルの初期化、チャネルのクローズ
* チャネル利用者では値の読み込みだけする

### p90
* スライスの範囲を分割することで mutex を使わずに済む

### p143
* ctx.Deadline() を使うことでタイムアウト発生前に処理を中断できる

### p151
* エラーの型を確認することで、そのエラーが想定内のエラーなのかバグなのかを区別する
