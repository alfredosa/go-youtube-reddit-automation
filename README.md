# E2E Content Creation App

## Description

This project is a fully developed bot that extracts data from Reddit, the text is then used to get text-to-speech, and Images, and together with Video footage it creates a ready Video to be uploaded to Youtube and Instagram Automatically. No human interaction.

![Video](example/resultwsound.mp4)

Through a Postgresql, we track Reddit Posts. This allows us to see and use data for sentiment analysis down the line, and to be able to never upload Posts twice.

## Features

- Feature 1: Pull all news of the day
- Feature 2: Filter the news and make a video with max 60 seconds of sound
- Feature 3: Create a video snippet of the exact duration of each news + transition sound
- Feature 4: Create a Banner Breaking News - With title for each news
- Feature 5: Pull all images from the related news (2 if possible) and create a video snippet with the images
- Feature 4: Add the Banner to the video snippet
- Feature 6: Concatenate all videos with an intro and sound
- Feature 7: Create a video with all news of the day

![Video Example](example/resultwsound.mp4)

## Dependencies

- FFMPEG (for video and audio)
- youtubedr (for video download) golang version

## Installation

Instructions on how to install and set up your project.

## Usage

Instructions on how to use your project. Include examples.

## TODOs

- [ ] Add a `LICENSE` file
- [x] Add a screenshot taker from reddit
- [x] Add a video downloader from youtube
- [x] DB connection to start saving posts on go-live and have info for data science
- [x] Add a `Dockerfile` to build and run the project WILL NOT DO
- [x] Add a `docker-compose.yml` to build and run the project WITH db

## Contributing

Information about how others can contribute to your project.

## License

Information about the license.
