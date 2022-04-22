// Package customerr provides custom error type and utility functions for error handling.
//
// Examples:
//
// 	func innerFunc(n string) error {
//		err := errors.New("inner error")
// 		if n == "fatal" {
// 			return New(err, []Tag{"fatal"}, "fatal error")
// 		}
// 		if n == "retryable" {
// 			return New(err, []Tag{"retryable"}, "retryable error")
// 		}
//		return nil
// 	}
//
// 	func outerFunc(n string) error {
// 		err := innerFunc(n)
// 		if err != nil {
// 			if HasTag(err, "retryable") {
// 				return New(err, nil, "invalid input")
// 			}
// 			return New(err, nil, "unknown error")
//		}
//		return nil
// 	}
//
// 	func main() {
// 		err := innerFunc(os.Args[1])
// 		if err != nil {
// 			if HasTag(err, "fatal") {
// 				log.Fatalln(err)
//				return
//			}
//			if HasTag(err, "retryable") {
//				fmt.Println("Please retry")
//				return
//			}
//			fmt.Println(err)
//		}
// 	}
// 	func handleError()
package customerr
