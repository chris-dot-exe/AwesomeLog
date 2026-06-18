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
	ShowTimestamp(false)
	bar := Foo{
		Foo: "Test",
		Bar: "Test",
		Foobar: struct {
			Meeps []string
		}{[]string{"Meep", "Meep2", "Meep2.1"}},
	}

	PrettyPrint(INFO, bar)
	// Output:
	// [INFO] [
	//   {
	//     "Foo": "Test",
	//     "Bar": "Test",
	//     "Foobar": {
	//       "Meeps": [
	//         "Meep",
	//         "Meep2",
	//         "Meep2.1"
	//       ]
	//     }
	//   }
	// ]
}

func ExamplePrint() {
	ShowTimestamp(false)
	Print(INFO, "hello ")
	Print(INFO, "world!")
	// Output:
	// [INFO] hello [INFO] world!
}

func ExamplePrint_withoutTimestamp() {
	ShowTimestamp(false)
	Print(INFO, "hello ")
	Print(INFO, "world!")
	// Output:
	// [INFO] hello [INFO] world!
}

func ExamplePrintln() {
	ShowTimestamp(false)
	Println(INFO, "line 1")
	Println(INFO, "line 2")
	// Output:
	// [INFO] line 1
	// [INFO] line 2
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
	ShowTimestamp(false)
	ShowCaller(false)
	Printf(DEBUG, "Hello %s!\n", "world")
	// Output:
	// [DEBUG] Hello world!
}

func ExampleSetTimeFormat() {
	ShowTimestamp(true)
	ShowCaller(false)
	SetTimeFormat("2006/01/02 15:04:05.000000")
	// Cannot easily test exact time output, so we mock time or skip output block
	// We just ensure the function works.
	Printf(DEBUG, "Hello %s!\n", "world")
}

func ExampleSetLevelConfig() {
	lvlConfig := DefaultLevelConfig()

	lvlConfig.Debug.ShowLineNumber = false
	lvlConfig.Debug.ShowFunctionName = true
	lvlConfig.Debug.ShowFilePath = false

	SetLevelConfig(lvlConfig)
}
