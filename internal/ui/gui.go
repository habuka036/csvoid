package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"csvoid/internal/exporter"
	"csvoid/internal/jsonflatten"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func toExporterTableRows(rows []jsonflatten.TableRow) []exporter.TableRow {
	out := make([]exporter.TableRow, len(rows))
	for i, r := range rows {
		out[i] = exporter.TableRow(r)
	}
	return out
}

func Launch() {
	myApp := app.New()
	myWindow := myApp.NewWindow("JSON to CSV/Excel Converter")

	// ファイルパス表示欄（複数ファイル用）
	selectedFiles := []string{}
	fileList := widget.NewMultiLineEntry()
	fileList.SetPlaceHolder("選択されたJSONファイルのパス...")

	// ファイル選択ボタン
	openFileBtn := widget.NewButton("ファイルを開く", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				path := reader.URI().Path()
				if !contains(selectedFiles, path) {
					selectedFiles = append(selectedFiles, path)
					fileList.SetText(strings.Join(selectedFiles, "\n"))
				}
			}
		}, myWindow)
	})

	// ドラッグ&ドロップ案内ラベル
	dropArea := widget.NewLabel("(ウィンドウにJSONファイルをドラッグ＆ドロップしてください)")

	// ドラッグ&ドロップでファイル取得（ウィンドウ全体で受け付ける）
	myWindow.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		for _, u := range uris {
			if strings.HasSuffix(strings.ToLower(u.Path()), ".json") && !contains(selectedFiles, u.Path()) {
				selectedFiles = append(selectedFiles, u.Path())
			}
		}
		fileList.SetText(strings.Join(selectedFiles, "\n"))
	})

	// 出力形式ラジオボタン
	formatRadio := widget.NewRadioGroup([]string{"CSV", "EXCEL"}, func(selected string) {})
	formatRadio.Selected = "CSV"

	// 出力先フォルダ
	outputPath := widget.NewEntry()
	outputPath.SetPlaceHolder("出力先のフォルダパス...")
	selectOutBtn := widget.NewButton("出力先フォルダ選択", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				outputPath.SetText(uri.Path())
			}
		}, myWindow)
	})

	// ステータス表示
	statusLabel := widget.NewLabel("")

	// 変換実行ボタン
	convertBtn := widget.NewButton("変換を実行", func() {
		statusLabel.SetText("")
		if len(selectedFiles) == 0 {
			dialog.ShowInformation("エラー", "ファイルを選択してください", myWindow)
			return
		}
		if outputPath.Text == "" {
			dialog.ShowInformation("エラー", "出力先フォルダを指定してください", myWindow)
			return
		}

		format := formatRadio.Selected // "CSV" or "EXCEL"
		successCount := 0
		failCount := 0
		for _, jsonFile := range selectedFiles {
			// ファイル読み込み
			data, err := os.ReadFile(jsonFile)
			if err != nil {
				dialog.ShowError(err, myWindow)
				failCount++
				continue
			}
			// JSONパース→フラット化
			var obj interface{}
			if err := json.Unmarshal(data, &obj); err != nil {
				dialog.ShowError(err, myWindow)
				failCount++
				continue
			}
			rows := jsonflatten.FlattenTable(obj)
			expRows := toExporterTableRows(rows) // convert to exporter.TableRow
			if err != nil {
				dialog.ShowError(err, myWindow)
				failCount++
				continue
			}
			// 出力ファイル名
			base := filepath.Base(jsonFile)
			base = strings.TrimSuffix(base, filepath.Ext(base))
			var outPath string
			if format == "CSV" {
				outPath = filepath.Join(outputPath.Text, base+".csv")
				f, err := os.Create(outPath)
				if err != nil {
					dialog.ShowError(err, myWindow)
					failCount++
					continue
				}
				if err := exporter.ExportCSV(expRows, f); err != nil {
					dialog.ShowError(err, myWindow)
					f.Close()
					failCount++
					continue
				}
				f.Close()
			} else {
				outPath = filepath.Join(outputPath.Text, base+".xlsx")
				f, err := os.Create(outPath)
				if err != nil {
					dialog.ShowError(err, myWindow)
					failCount++
					continue
				}
				if err := exporter.ExportExcel(expRows, f); err != nil {
					dialog.ShowError(err, myWindow)
					f.Close()
					failCount++
					continue
				}
				f.Close()
			}
			successCount++
		}
		result := "変換完了: " + toString(successCount) + "件成功"
		if failCount > 0 {
			result += " / " + toString(failCount) + "件失敗"
		}
		statusLabel.SetText(result)
		dialog.ShowInformation("処理結果", result, myWindow)
	})

	// ファイルリストと「開く」ボタンをフォームの1行にする
	fileSelectRow := container.NewBorder(
		nil, nil, nil, openFileBtn, fileList,
	)
	// 出力先欄も同様にボタンとEntryをまとめる
	outputRow := container.NewBorder(
		nil, nil, nil, selectOutBtn, outputPath,
	)

	form := widget.NewForm(
		widget.NewFormItem("ファイル一覧", fileSelectRow),
		widget.NewFormItem("出力形式", formatRadio),
		widget.NewFormItem("出力先フォルダ", outputRow),
	)

	content := container.NewVBox(
		dropArea,
		form,
		convertBtn,
		statusLabel,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}

// ユーティリティ：スライスに含まれるか
func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// int→string
func toString(i int) string {
	return strconv.Itoa(i)
}
