package log

type Foo struct {
	Foo    string
	Bar    string
	Foobar struct {
		Meeps []string
	}
}

func ExampleSetDefaultLevel() {
	SetDefaultLevel(INFO)
}

func ExampleSetLogLevel() {
	SetLogLevel(WARN)
}

func ExampleSetLogLevelByString() {
	SetLogLevelByString("WARN")
}

func ExamplePrettyPrint() {
	bar := Foo{
		Foo: "Test",
		Bar: "Test",
		Foobar: struct {
			Meeps []string
		}{[]string{"Meep", "Meep2", "Meep2.1"}},
	}

	PrettyPrint(INFO, bar)
	// Output:
	// 2022/02/15 22:32:38 [INFO] [
	//  {
	//    "Foo": "Test",
	//    "Bar": "Test",
	//    "Foobar": {
	//      "Meeps": [
	//        "Meep",
	//        "Meep2",
	//        "Meep2.1"
	//      ]
	//    }
	//  }
	//]
}

func ExamplePrint() {
	Print(INFO, "hello ")
	Print(INFO, "world!")
	// Output:
	// 2022/02/15 22:46:51 [INFO] hello
	// 2022/02/15 22:46:51 [INFO] world!
}

func ExamplePrint_withoutTimestamp() {
	ShowTimestamp(false)
	Print(INFO, "hello ")
	Print(INFO, "world!")
	// Output:
	// [INFO] hello [INFO] world!
}

func ExamplePrintln() {
	Println(INFO, "line 1")
	Println(INFO, "line 2")
	// Output:
	// 2022/02/15 22:46:51 [INFO] line 1
	// 2022/02/15 22:46:51 [INFO] line 2
}

func ExamplePrintln_withoutTimestamp() {
	ShowTimestamp(false)
	Println(INFO, "line 1")
	Println(INFO, "line 2")
	// Output:
	// [INFO] line 1
	// [INFO] line 2
}

func ExamplePrintf() {
	Printf(DEBUG, "Hello %s!\n", "world")
	// Output: 2022/02/15 23:18:02 [DEBUG][main.go:main:47] Hello world!
}

func ExampleSetLevelConfig() {
	lvlConfig := DefaultLevelConfig()

	lvlConfig.Debug.ShowLineNumber = false
	lvlConfig.Debug.ShowFunctionName = true
	lvlConfig.Debug.ShowFilePath = false

	SetLevelConfig(lvlConfig)
}
