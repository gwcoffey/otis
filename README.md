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
author: 
  name: Wendy Writer
  surname: Writer

address: |-
  1234 Sesame Street
  My Town, Arizona 90210
  555-555-1212
  me@example.com
```

This content is used on the title page of compiled manuscripts.

## The Manuscript

The source manuscript in an otis project is just a collection of content files, metadata files, and folders. In its simplest form a manuscript can be just this:

* `/manuscript`
  * `/work.yml` defines the work metadata (title, etc…) 
  * `/00-content.md` the actual content

> Note: Notice the content file starts with `00`. This is essential in Otis and ensures that books, chapters, scenes, etc… are ordered correctly when editing and when compiled. Otis helps you manage these numbers. 

Of course manuscripts often benefit from more structure. An otis manuscript has four basic concepts (each explained more fully below):

* Works: independently published works within the same manuscript, like books in a series.
* Folders: for structuring your manuscript, folders are for *you* and have no bearing on the compiled output
* Chapters: for structuring the compiled output
* Scenes: actual content

### Manuscript Works

A manuscript can contain multiple *works* (think multiple books in a series):

* `/manuscript/00-my-book`
    * `/work.yml`
    * `00-content.md`
* `/manuscript/01-the-second-book`
    * `/work.yml`
    * `00-content.md`
* `/manuscript/02-the-final-book`
    * `/work.yml`
    * `00-content.md`

Works have a title and are compiled into distinct manuscript files. For example, the above structure, when compiled, would produce three separate documents.

The `work.yml` looks like this:

```yaml
title: My New Book
runningTitle: Book
author: Wendy Writer
authorSurname: Writer
```

It has these properties:
* `title` the full title of the work
* `runningTitle` (optional) the short title, displayed on top of each page in the compiled output
* `author` (optional) the pen name of the author (overrides the author name from the otis config)
* `authorSurname` (optional) the surname of the author, displayed on top of each page in the compiled output (overrides the surname from the otis config)

The `/manuscript` directory can *either* have a single `work.yml` at the top level, or one or more folders at the top level, each of which has a `work.yml`.

### Manuscript Folders

In addition to folders per work, you can use folders to organize content *within* a work. These folders have no impact on the structure of the compiled output. They are useful to *you* to organize, outline, and structure your writing environment.

Here's an example:

* `/manuscript`
    * `/work.yml`
    * `/00-act-1`
        * `/00-content.md`
    * `/00-act-2`
        * `/00-content.md`

Like works, folders are always numbered with a two-digit prefix. (If you need more than 100 you need to subdivide into more folders.) They are named with `lower-kebab-case`.

### Manuscript Scenes

Scenes are markdown files within works and folders. They are numbered and have the `.md` file extension. They may also have a name, in `lower-kebab-case`, to help you identify them. (This name has no bearing on the final output.)

Here's an example of a work with several scenes:

* `/manuscript`
    * `/work.yml`
    * `/00-wrap-her-in-package-of-lies.md`
    * `/01-this-is-not-love.md`
    * `/02-anna-begins-to-change-my-mind.md`

When compiled, scenes are separated by a *scene break*, which in the standard manuscript format is `#` centered on a line by itself.

> Note: While scenes are written in markdown, otis has very limited actual markdown support. When compiling to HTML otis uses a standard markdown processor. But when compiling to other formats, otis processes the markdown itself, and currrently only supports `*emphasis*` and `> blockquotes`. All other markdown will be copied to the output unchanged.

### Manuscript Chapters

Chapters are only loosely tied to the folder hierarchy. While you *may* structure your manuscript with one folder per chapter, this is not strictly required. Sometimes you may want to use folders to outline a complex work, and apply chapters over that structure, and otis support this.

Essentially you can think of the manuscript as a long list of scenes. When compiling output, otis walks through the folder hierarchy and the scenes within those folders. Chapters are waypoints along this walk.

You identify a chapter with a `chapter.yml` metadata file placed anywhere in the folder hierarchy. This declares that everything from the first scene in this folder, to the next chapter marker, belongs to this chapter.

Otis only has one rule with chapter markers, which is: *if you have any chapters, they must start with the first scene*. You cannot have scenes at the beginning of the manuscript before the first chapter.

Here's a simple example:

* `/manuscript`
    * `/work.yml`
    * `00-chapter-1`
        * `chapter.yml`
        * `00-content.md`
    * `01-chapter-2`
        * `chapter.yml`
        * `00-content.md`
    * `03-chapter-3`
        * `chapter.yml`
        * `00-content.md`

This manuscript has one work with three chapters. Each chapter is one folder with a `chapter.yml` and one or more scenes.

Here's a more complex example where the chapter markers don't align directly with the top-level folder hierarchy:

* `/manuscript`
    * `/work/yml`
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
$ o -h
```

The command expects to be run from within a project directory (which it identified by searching up the directory structure for the nearest `otis.yml`).

> Note: If you really want to you can use `otis` on a project you're not *in* using the `--project` command line switch.

### Initializing a New Project

To start a new project, create an empty directory, switch to it, and run:

```shell
$ otis init
```

This will create the necessary configuration for a single-work project with one scene. You can edit the metadata files as needed, and add scenes.

If you plan to create a multi-work project, use the `--works` option:

```shell
$ otis init --works 3
```

This version will create individual folders for each work with a work metadata file in each.

### Working with Scenes

```shell
# add a scene to the end
$ otis touch PATH NAME

# add a scene inserted
$ otis touch PATH NAME --at NEW_INDEX

# move a scene in the same folder
$ otis mv PATH --at NEW_INDEX

# move a scene to another folder
$ otis mv PATH PATH

# move a scene to another folder and insert
$ otis mv PATH PATH --at NEW_INDEX

# split a scene (insert "###" in scene file first)
$ otis split PATH

# combine multiple scenes into one
$ otis join PATH PATH...

# normalize all scene/folder numbers
$ otis normalize [--recursive] PATH

# insert a chapter
$ otis chapter PATH

# show outline
$ otis ls

# show table of contents
$ otis ls --chapter
```
 

### Working with Folders

TODO

### Working with Chapters

TODO

### Counting Words

Otis can count the words in your manuscript:

```shell
$ o wordcount
```

This will show the wordcount for every scene in the manuscript, along with subtotals at the folder level and totals at the work level. 

You can target a single work by naming it as an argument:

```shell
$ o wordcount 00-book-1
```

And you can get your counts by *chapter* instead of *scene*/*folder* with `--chapter`:

```shell
$ o wordcount --chapter
```

### Compiling

While some people (maybe just me) find *writing* in simple text files and using git for revision management, branching, etc… a breath of fresh air, these are not suitable formats for sharing your work with others. Otis can *compile* your manuscript into standard readable forms.

```shell
$ o compile
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
$ o compile --format RTF
```

Normally otis names the output file with the current date appended to the end. But you can change this with the `--tag` switch.

```shell
$ o compile --tag "draft1"
```

The output file name will still be based on the title of the work, but it will have the tag name appended instead of the date. 