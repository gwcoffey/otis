# Otis

Otis is a simple and opinionated long-form writing tool meant to serve the 0.0000000002% of the population that likes to write long-form text (book manuscripts, short stories, etc…) and prefers using simple text files + git for the actual writing.

Otis is a work in progress, and built for my own use. It was born out of my frustration with traditional writing tools like Word and Scrivener. These are very powerful in some ways (many of which I don't care much about). And they are very weak when it comes to revision tracking. Since I'm very comfortable in git, I decided to try doing my writing in git-friendly formats instead. I prefer to work this way, and otis is my tool to make this practicable. 

> Note: This is a work in progress, not super well-tested, and may evolve in breaking ways over time.

## Otis Projects

An otis project is just an ordinary directory with these special files and folders:

* `otis.yml` the otis configuration
* `/manuscript` the actual manuscript content
* `/dist` where compiled output goes

The project can contain other files and folders. Otis will ignore them.

## Configuring the Project

An `otis.yml` file at the root of the project configures otis and the manuscript. It looks like this:

```yml
title: Title of Mine
runningTitle: Mine
author: 
  name: Wendy Writer
  surname: Writer

address: |-
  1234 Sesame Street
  My Town, Arizona 90210
  555-555-1212
  me@example.com
```

This content is used on the title page and headers of compiled manuscripts.

> Tip: If the first line of the `otis.yml` file is a comment with the form `#import path/to/a/file` then the referenced file will be imported into the otis.yml. This is handy when you want to share author information across several projects, as in a multi-book project or a multi-project repo.

## The Manuscript

The source manuscript in an otis project is just a collection of content files, metadata files, and folders. In its simplest form a manuscript can be just this:

* `/manuscript`
  * `/00-content.md` the actual content

> Note: Notice the content file starts with `00`. This is essential in Otis and ensures that books, chapters, scenes, etc… are ordered correctly when editing and when compiled. Otis helps you manage these numbers. 

Of course manuscripts often benefit from more structure. An otis manuscript has four basic concepts (each explained more fully below):

* Folders: for structuring your manuscript, folders are for *you* and have no bearing on the compiled output
* Chapters: for structuring the compiled output
* Scenes: actual content

### Manuscript Folders

You can use folders to organize content however you like. These folders have no impact on the structure of the compiled output. They are useful to *you* to organize, outline, and structure your writing environment.

Here's an example:

* `/manuscript`
    * `/00-act-1`
        * `/00-content.md`
    * `/00-act-2`
        * `/00-content.md`

Folders are always numbered with a two-digit prefix. (If you need more than 100 you need to subdivide into more folders.) They are named with `lower-kebab-case`.

### Manuscript Scenes

Scenes are markdown files within the manuscript or its folders. They are numbered and have the `.md` file extension. They may also have a name, in `lower-kebab-case`, to help you identify them. (This name has no bearing on the final output, and is optional.)

Here's an example of a manuscript with several scenes:

* `/manuscript`
    * `/00-wrap-her-in-package-of-lies.md`
    * `/01-this-is-not-love.md`
    * `/02-anna-begins-to-change-my-mind.md`

When compiled, scenes are separated by a *scene break*, which in the standard manuscript format is `#` centered on a line by itself.

> Note: While scenes are written in markdown, otis has very limited actual markdown support. When compiling to HTML otis uses a standard markdown processor. But when compiling to other formats, otis processes the markdown itself, and currrently only supports `*emphasis*` and `> blockquotes`. All other markdown will be copied to the output unchanged.

### Manuscript Chapters

Chapters are only loosely tied to the folder hierarchy. While you *may* structure your manuscript with one folder per chapter, this is not strictly required. Sometimes you may want to use folders to outline a complex manuscript, and apply chapters over that structure, and otis support this.

Essentially you can think of the manuscript as a long list of scenes. When compiling output, otis walks through the folder hierarchy and the scenes within those folders. Chapters are waypoints along this walk.

You identify a chapter with a `chapter.yml` metadata file placed anywhere in the folder hierarchy. This declares that everything from the first scene in this folder, to the next chapter marker, belongs to this chapter.

Otis only has one rule with chapter markers, which is: *if you have any chapters, they must start with the first scene*. You cannot have scenes at the beginning of the manuscript before the first chapter.

Here's a simple example:

* `/manuscript`
    * `00-chapter-1`
        * `chapter.yml`
        * `00-content.md`
    * `01-chapter-2`
        * `chapter.yml`
        * `00-content.md`
    * `03-chapter-3`
        * `chapter.yml`
        * `00-content.md`

This manuscript has three chapters. Each chapter is one folder with a `chapter.yml` and one or more scenes.

Here's a more complex example where the chapter markers don't align directly with the top-level folder hierarchy:

* `/manuscript`
    * `00-act-1`
        * `chapter.yml`
        * `00-introduction`
            * `00-content.md`
    * `01-act-2`
        * `00-conflict`
            * `chapter.yml` 
            * `00-content.md`
        * `01-defeat`
            * `chapter.yml`
            * `00-content.md`
        * `02-redemption`
            * `chapter.yml`
            * `00-content.md`
    * `02-act-3`
        * `chapter.yml`
        * `00-content.md`

This manuscript has 5 chapters, three of which are within act 2. When you use chapters this way, you can structure your manuscript around whatever outlining style you prefer, and insert chapter boundaries later, or move them, without adjusting the outline structure.

A `chapter.yml` looks like this:

```yml
title: My Chapter Title
numbered: false
```

It has these properties:
* `title` the title of the chapter, displayed at the top fo the first page of the chapter
* `numbered` (optional, defaults to `true`) when true, otis will display, eg, `Chapter 1` above the chapter title. 

> Note: You can mix numbered and un-numbered chapters. For instance, you may have an unnumbered "Epilogue", "Introduction", etc…, then a series of numbered chapters, and then an unnumbered "Afterword". 

## The `otis` Command

The otis command line tool helps you work with your manuscript. You can:

* initialize an otis project
* create, move, and split scenes
* add folders, chapters, etc…
* count words by folder and scene, or by chapter
* compile the manuscript into a readable or submittable form

All these are available via the `otis` command. To get started, ask otis for help:

```shell
$ otis -h
```

### Initializing a New Project

To start a new project, create an empty directory, switch to it, and run:

```shell
$ otis init
```

This will create the necessary configuration for a manuscript with one scene. You can edit the metadata files as needed, and add scenes.

### Creating Scenes

You can add new scenes with the `touch` command:

```shell
$ otis touch manuscript/ "my new scene"
$
```

This will add a new file at the root of the manuscript with a name like `00-my-new-scene.md`. The scene will, by default, be placed at the end of the target folder. In other words, the index number will be the largest current index number plus one.

You can specify an index at which to insert the scene instead:

```shell
$ otis touch --at 3 manuscript/ "my new scene"

About to change:

  RENAME 03-an-existing-scene.md → 04-an-existing-scene.md
  RENAME 04-another-scene.md → 05-another-scene.md
     ADD 03-my-new-scene.md

OK to proceed? [Y/n]: 
```
> Note: As shown above, any time otis will make more than one change to the manuscript to accommodate
> a command, it will first show the list of changes it is about to make and prompt you for
> confirmation. You can use the `--force` switch to bypass this prompt if you prefer. 

This time, the scene index will be `03` and any existing scenes will be adjusted forward as necessary.

### Creating Folders

You can add new folders with the `mkdir` command:

```shell
$ otis mkdir manuscript/ "my new folder"
$  
```

This will create a new folder in the target directory named something like `00-my-new-folder`. Again, Otis will number the folder so it is added to the *end* of the target. And again you can use `--at` to insert it somewhere else:

```shell
$ otis mkdir --at 0 manuscript/ "my new folder"
About to change:

  RENAME 00-my-first-scene.md → 01-my-first-scene.md
  RENAME 01-my-second-scene.md → 02-my-second-scene.md
     ADD 00-my-new-folder

OK to proceed? [Y/n]: 
```

### Moving Scenes and Folders

You can move scenes and folders using the `mv` command. You can use `mv` to **renumber** something:

```shell
$ otis mv --at 5 manuscript/04-my-scene.md
```

This will rename the specified scene or folder with the new index number, and rename any other scenes/folders as necessary to make room. Since scene names minus the index may not be unique, this sometimes requires moving a file to a temporary location, making space, and moving it back. Again, Otis will show you everything it is going to do before it does it.

You can also move a scene to another folder:

```shell
$ otis mv manuscript/00-act-1/00-my-first-scene.md manuscript/01-act-2
```

This will move the scene out of its current folder (renumbering things as needed) and add it to the *end* of the target folder.

Finally you can combine `--at` with a target folder to control the final index of the scene:

```shell
$ otis mv --at 0 manuscript/00-act-1/00-my-first-scene.md manuscript/01-act-2 
```

> Note: `mv` works the same with both scenes and folders. 

### Counting Words

Otis can count the words in your manuscript:

```shell
$ otis wordcount
```

This will show the wordcount for every scene in the manuscript, along with subtotals at the folder level.

You can count words by *chapter* instead of *scene*/*folder* with `--chapter`:

```shell
$ otis wordcount --chapter
```

### Compiling

While some people (maybe just me) find *writing* in simple text files and using git for revision management, branching, etc… a breath of fresh air, these are not suitable formats for sharing your work with others. Otis can *compile* your manuscript into standard readable forms.

```shell
$ otis compile
```

This puts the compiled file in `/dist` within the project folder. By default, it is named `{title}_{date}.{format}`.  

It supports four output formats:

* `PDF` (default) a pdf in standard manuscript format
* `HTML` a single page HTML file that mimics standard manuscript format
* `RTF` a rich text file that can be opened in Microsoft Word, LibreOffice, Pages for macOS, etc… It follows standard manuscript format.
* `TEX` a LaTeX document (this is the source file for the PDF version, but you can produce this directly if it is useful to you).

> PDF Note: The pdf output is produced by compiling to LaTeX and then processing the `.tex` file with the `pdflatex` command. You need a LaTeX installation (with `pdflatex` in your path) for this to work.

> RTF Note: When opening in Pages for macOS, the RTF format does not currently include page headers. I haven't been able to find a way to make this work.

You specify the format you want with the `--format` option:

```shell
$ otis compile --format RTF
```

Normally otis names the output file with the current date appended to the end. But you can change this with the `--tag` switch.

```shell
$ otis compile --tag "draft1"
```

The output file name will still be based on the title of the manuscript, but it will have the tag name appended instead of the date. 