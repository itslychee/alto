# cmd/alto's scope

The command line program `alto` consists of numerous functions and variables to help
aid the user in organizing their music. The original purpose of this project was to provide a simple, yet 
powerful language for organizing music at a breeze.

## Variables

Along with the default variables provided in the `dsl` package, alto provides the following:

* `%trackcurrent%` — current audio track number
* `%tracktotal%` — the total amount of tracks in the disc
* `%disccurrent%` — current disc number
* `%disctotal%` — total amount of discs
* `%year%` — year this work was produced
* `%comment%` — metadata comment
* `%format%` — audio container name
* `%composer%` — composer of this work
* `%genre%` — type of genre of this work
* `%albumartist%` — albumartist of this work
* `%album%` — album name of this work
* `%artist%` — artist of this work
* `%title%` — title of this work
* `%filetype%` — Type of file of this work, lowercased by default
* `%alto_dest%` — Destination directory alto is writing to
* `%alto_source%` — Source directory alto is reading from

Please note that the definitions of these variables are what are *expected*, that doesn't mean you will not
receive the values you expect. This is why alto was created, to handle arbitrary metadata with relative ease.

## Functions

### \<skip>
Skip automatically ceases execution and proceeds to the next file

***Returns:*** `ErrSkip`

### \<print **{arg...}**>
Print writes to the stderr, mainly useful for debugging path constructs

- {arg...} Varadic sequence of arguments that will be joined together with `""`

***Returns:*** Nothing

### \<clean **{arg}**>
The clean function is mainly for illegal character sanitization, to help mitigate 
using arbitrary metadata without any file creation errors.

- {arg} String that the function will clean

***Returns:*** {arg} stripped of any illegal filename characters, if any

### \<exists **{path}**>
The exists function checks whether or not `{path}` exists

- {path} Filepath

***Returns:*** Returns `{path}` if it exists, otherwise it returns an empty string

### \<uniqueFp **{path}**>

uniqueFp is a function that returns a non-existent filepath, however there's more to it then what meets
the eye. It copies the current scope, which is what contains variables and functions, and sets a variable 
`%index%` to the total times **{path}** has been executed. 

This function is quite unideal for setting variables that retain outside it, but perhaps it may
be useful for preventing the scope from getting polluted if you require a better contextual state.

- **{path}** A node, typically a group, that will be executed until it returns a filepath that doesn't exist

***Returns:*** A filepath that doesn't exist