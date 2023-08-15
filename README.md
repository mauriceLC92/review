# Review

[![Go Reference](https://pkg.go.dev/badge/github.com/mauriceLC92/review.svg)](https://pkg.go.dev/github.com/mauriceLC92/review) [![Go Report Card](https://goreportcard.com/badge/github.com/mauriceLC92/review)](https://goreportcard.com/report/github.com/mauriceLC92/review)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mauriceLC92/review
)

The `review` package provides an easy and intuitive way to do monthly reviews. Easily track how your month went and reflect on past months to see if you are progressing.

## Why
Conducting a monthly review is an essential practice for tracking progress, reflecting on achievements, and setting new goals, both personally and professionally.

However, the process can be time-consuming and disorganized. The `review` package simplifies this practice by providing a streamlined CLI to create, manage, and review monthly assessments all in one place.

Whether you're an individual seeking to improve personal growth or a team leader focusing on collective advancement, this package offers a structured and efficient way to make monthly reviews a productive and insightful routine.
## Installation
```sh
go install github.com/mauriceLC92/review/cmd/review
```

## Usage

Start a new review if one has not been done, check when the next date is due or initiate your next review.
```bash
review
```

Example output:
```
You have not done a review yet! Let's get you started

How are you today?
Doing well, but the month has been super busy!

What was your biggest win this month?
Sticking to my daily hour of writing Go.

Thanks for the review! See you on 2023-09-05 11:23:09.757929 +0200 SAST for the next one!
```
```
You review is not due until 2023-08-13 00:00:00 +0000 UTC. See you then! 
```

List previous reviews which have been filled out.
```bash
review list
```
Example output:
```
Date: 13-07-2023
Questions:
Title: How are you today?
Answer: Great!
----------------------------------------------------
Title: What was your biggest win this month?
Answer: 1 hour of Go each day
----------------------------------------------------
```

Display the CLI commands available to you.
```bash
review help
```
```
review - Initiate a new review or check when the next is due.
list - List previous reviews which have been filled out.
help - Display the CLI commands available to you.
```
## License

This package is licensed under the MIT License - see the LICENSE file for details.