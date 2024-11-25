# Go Multithreaded File Copy Tool

## Description
This project is a Go-based multithreaded tool designed to copy files from a source directory to a destination directory. It uses parallelism to speed up the copying process and ensures robustness through retry mechanisms. The program is configured through environment variables, supports logging, and displays the progress of the copying process.

## Features
- **Multithreading**: Uses multiple workers (goroutines) to perform file copies concurrently, enhancing speed.
- **Retry Mechanism**: Each file is attempted up to three times in case of failure.
- **Logging**: Logs all operations, including successful copies and errors, to a log file.
- **Progress Tracking**: Displays the progress of file copying, including estimated remaining time.

## Requirements
- Go 1.16 or newer
- Git
- A `.env` file to specify the necessary environment variables

## Setup Instructions

### Step 1: Clone the Repository
First, clone the repository to your local machine:
```sh
$ git clone https://github.com/darskip/go-multithreaded-copy.git
$ cd go-multithreaded-copy
```

### Step 2: Install Dependencies
Install the `godotenv` package used to load environment variables:
```sh
go get github.com/joho/godotenv
```

### Step 3: Create a `.env` File
Create a `.env` file in the root directory of the project to specify the configuration:
```env
SOURCE_DIR=P:\lossless\
DEST_DIR=\\192.168.133.230\Sony\Sources Demat\Sony.ddex\lossless\
FILES_LIST_PATH=sony2024.txt
```
- `SOURCE_DIR`: The source directory containing the files to be copied.
- `DEST_DIR`: The destination directory where the files will be copied.
- `FILES_LIST_PATH`: The path to the file containing a list of files to be copied.

### Step 4: Run the Program
To run the program, use the following command:
```sh
go run main.go
```

### Step 5: Build the Program (Optional)
If you want to build the project into an executable:
```sh
go build -o file-copy-tool
```
This will create an executable named `file-copy-tool` that you can use to run the program.

## Project Structure
- **main.go**: Contains the main logic for multithreaded copying.
- **.env**: Environment variables to configure source, destination, and list paths.
- **copy.log**: Log file to track the progress and errors during file copying.

## Usage
The tool will read the list of files from the file specified in `FILES_LIST_PATH`, then use a pool of workers to copy these files from `SOURCE_DIR` to `DEST_DIR`. Any failures are retried up to a specified number of times (default is 3). The progress and any errors encountered will be logged to the console and a log file (`copy.log`).

## Notes
- Ensure that the `SOURCE_DIR` and `DEST_DIR` are accessible from the system where the program is run.
- If the `DEST_DIR` is a network path, proper permissions are required to access the network share.
- For long paths on Windows, you might need to use UNC paths prefixed with `\\?\`.

## License
This project is open source and available under the [MIT License](LICENSE).

## Contributing
Feel free to submit issues or pull requests for any improvements or additional features you would like to see.

## Author
- **darskip**: [GitHub Profile](https://github.com/darskip)

