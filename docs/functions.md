# DSL Functions

This page will talk about the predefined functions at `dsl.DefaultFunctions`. Of course, you are entitled
to using a different implementation if you so wish to.

## \<eq  **{cond1}** **{cond2}** **{body}**>
## \<neq **{cond1}** **{cond2}** **{body}**>
## \<gt  **{cond1}** **{cond2}** **{body}**>
## \<lt  **{cond1}** **{cond2}** **{body}**>
## \<gte **{cond1}** **{cond2}** **{body}**>
## \<lte **{cond1}** **{cond2}** **{body}**>

- \<eq>  Equal To
- \<neq> Not Equal To
- \<gt>  Greater Than
- \<lt>  Less Than
- \<gte> Greater or Equal To
- \<lte> Less or Equal To


These functions are all based on a common implementation, with a few exceptions of course. If the two
arguments meet the criteria set by the function, then execution will be passed onto **{body}**, otherwise
the selected function will return an empty result.

Exceptions:
* Only numeric values are allowed in comparisons that involve a range of numbers, e.g. gt, lt, etc

- **{cond1}** The first condition
- **{cond2}** The second condition
- **{body}** The body that which will be executed if the conditional succeeds

***Returns:*** ***{body}*** or an empty result



## \<fset **{name}** **{value}**>

The same implementation as `<set>` below, but it allows overwriting an already existing variable.

## \<set **{name}** **{value}**>

Set a variable in the scope, if it exists then an error will be returned.

- **{name}** Name of the new variable
- **{value}** Value of the new variable

***Returns:*** An empty result or an error if the variable exists


## \<trim **{value}**>

Trim removes any excess whitespace, the function is just a pipeline to Go's [`strings#TrimSpace`](https://pkg.go.dev/strings#TrimSpace)

- **{value}** String to be trimmed

***Returns:*** A clean result **{value}**, if it does not need to be trimmed then the value will be passed back

## \<exit>

Exit is a function that signals the program to quit execution

***Returns:*** `ErrExit`

## \<must **{...nodes}**>

The must function requires each node to return a perfect response, otherwise return an error stating which function
returned an imperfect response

A perfect response is a response not having an error but having a non-empty string

***Returns:*** **{...nodes}** joined by `""`

