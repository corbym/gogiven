Firstly, thanks for wanting to contribute!

There are only a few rules to helping out:

* Please write your tests first. And yes, I will be able to tell :) 
* Any code you submit without test coverage at at least 95% (matching the coverage overall) will not be approved! 
* Running go fmt is also a good idea.
* Public funcs and methods should always have godoc associated with them, starting with the name of the func/method. e.g. //MyFunc does so and so..
* Tests in the project do not necessarily match idiomatic go tests organisation, as I tend to write integration type unit tests rather than single package item unit tests. However, 
any changes to existing code should have unit tests added at the package item level (new structs and funcs) where possible.
* Any broken tests should be corrected before submitting.
* Keep code modular and small. Functions over ten lines will not be accepted.
* Follow [idiomatic go coding guildelines](https://golang.org/doc/effective_go.html) whereever possible

Most importantly, however, *have fun*!
