# gphotos-cli: A command line interface for Google Photos

A command line interface app for Google photos for doing bulk tasks such as mass filtered downloads
and mass uploads.

## Features

* Parallelized transfers
* Filtering by date and category type for downloads
* Album creation
* Uploads to designated albums (can only upload to albums created by this app due Google's API limitations)
* Exponential backoff for large uploads and downloads

# Installation

With Go already installed and setup, just call:

`go get -u http://github.com/trinhdrew1418/gphotos-cli`

Before usage, you first need to grant the application account authorization:

`gphotos-cli init`

# Usage
## `init`

This initializes the authorization token needed to access your account

`gphotos-cli init`

## `push`

Command to upload files to your google photos library.

`gphotos-cli push [-OPTIONS] [FILE 1] [FOLDER 1] ...`


Here are the available options:

* `-r` or `--recursive`

Recursively traverses folders in arguments for file uploads

* `-s` or `--select`

Pulls up a menu to select available albums from to upload to. NOTE: you can ONLY upload to
albums that you've created through this app. This is a limitation of the google photos api.

* `-v` or `--verbose`

This lists out all of the files you'll be uploading. Useful if you want to know which files have
valid extensions.

## `pull`

Command to download files

`gphotos-cli pull [-OPTIONS]`

Follow the prompts to select your filters. It will download your files into folders separated by year and month. Each
file will be named with its day of creation and time.

Here are the available options:

* `-d [PATH]`

define the directory path you want to download your files to.

## `albums`

Create albums by calling

`gphotos-cli albums create`

More subcommand coming later.

# Caveats

* You may unexpectedly hit a quota limit due to the application coming with some default credentials.
* Due to a google photos api limitation, you can only upload to albums created by the app
* When downloading files, its undetermined how many photos in total you`ll be downloading. The API uses pagination for large
requests and instead of waiting for possibly several page requests for the full download request, it is faster to do
concurrent requests for pages whilst downloading whats available.

#### TODOs

* More convenient credential replacement using `init`
* Mass moving existing photos to albums
* Search by metadata: specific filetypes, camera types, etc.
* Add options to change the max amount of parallel threads
