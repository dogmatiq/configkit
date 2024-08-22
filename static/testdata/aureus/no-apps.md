# No applications

This test that the static analyzer does not fail when the analyzed code does not
contain any Dogma applications.

```go au:input
package app

type NonApp struct{}
```

```au:output
(no applications found)
```
