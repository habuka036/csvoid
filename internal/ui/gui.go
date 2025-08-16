package ui

import (
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

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
                selectedFiles = append(selectedFiles, reader.URI().Path())
                fileList.SetText(strings.Join(selectedFiles, "\n"))
            }
        }, myWindow)
    })

    // ドラッグ&ドロップ案内ラベル
    dropArea := widget.NewLabel("(ウィンドウにJSONファイルをドラッグ＆ドロップしてください)")

    // ドラッグ&ドロップでファイル取得（ウィンドウ全体で受け付ける）
    myWindow.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
        for _, u := range uris {
            if strings.HasSuffix(u.Path(), ".json") {
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

    // 変換実行ボタン
    convertBtn := widget.NewButton("変換を実行", func() {
        // ここで変換処理を呼び出す
    })

    // fileListと「開く」ボタンをフォームの1行にする
    fileSelectRow := container.NewBorder(
        nil, nil, nil, openFileBtn, fileList,
    )

    // 出力先欄も同様にボタンとEntryをまとめる
    outputRow := container.NewBorder(
        nil, nil, nil, selectOutBtn, outputPath,
    )

    // FormでまとめるとEntry類が"ほどよい幅"に揃う
    form := widget.NewForm(
        widget.NewFormItem("ファイル一覧", fileSelectRow),
        widget.NewFormItem("出力形式", formatRadio),
        widget.NewFormItem("出力先フォルダ", outputRow),
    )

    // 全体レイアウト
    content := container.NewVBox(
        dropArea,
        form,
        convertBtn,
    )

    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(600, 400))
    myWindow.ShowAndRun()
}

/*
package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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
				selectedFiles = append(selectedFiles, reader.URI().Path())
				fileList.SetText(strings.Join(selectedFiles, "\n"))
			}
		}, myWindow)
	})

	// ドラッグ&ドロップ案内ラベル
	dropArea := widget.NewLabel("(ウィンドウにJSONファイルをドラッグ＆ドロップしてください)")

	// ドラッグ&ドロップでファイル取得（ウィンドウ全体で受け付ける）
	myWindow.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		for _, u := range uris {
			if strings.HasSuffix(u.Path(), ".json") {
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

	// 変換実行ボタン
	convertBtn := widget.NewButton("変換を実行", func() {
		// ここで変換処理を呼び出す
	})

	// レイアウト
	topRow := container.NewGridWithColumns(
		2,
		container.NewVScroll(fileList),
		openFileBtn,
	)
	outputRow := container.NewHBox(outputPath, selectOutBtn)
	widgetList := container.NewVBox(
		topRow,
		dropArea,
		formatRadio,
		outputRow,
		convertBtn,
	)

	myWindow.SetContent(widgetList)
	myWindow.Resize(fyne.NewSize(600, 350))
	myWindow.ShowAndRun()
}
*/