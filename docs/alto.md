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

### <uniqueFp {path}>

**`<uniqueFp ...>`** is a function that will call **`{path}`** forever until it returns a path
that doesn't exist. After the first iteration of calling **`{path}`** it will provide a variable 
called **`%index%`** to represent the count of iterations. Another unique thing about this
function is that it has a contained scope, which means that anything updated within **`{path}`** will
stay, so **`<fset ...>`** and **`<set ...>`** calls will not retain outside this function.

#### Example

```golang
// Path construct
<uniqueFp {%album%/%title%{ (%index%)}.%filetype%}>

// Output
Album/Title.flac      // %index% isn't present
Album/Title (1).flac  // %index% is present, set to 1
Album/Title (2).flac  // %index% is present, set to 2
...
```

### <exists {path}>

**`<exists {path}>`** is a function that will return **`{path}`** if it doesn't exists, otherwise
it will return an empty string.

#### Example

```golang
// Path construct
{<exists {%album%/%title%.%filetype%}>|<skip>}

// Output
Album/Title.flac // Function <skip> was not called, so this filepath doesn't exist yet
<skip> // Function <skip> was called
<skip> // Function <skip> was called
...
//