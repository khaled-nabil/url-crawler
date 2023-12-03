# url-crawler
Crawls a list of URLs weekly, saves them to file

## Usage
Create folder ~/mrcrowley with current user.
Folder can be changed in `commander.sh` file.

```bash
$ mkdir ~/mrcrowley
$ ./commander.sh
```

this uses default settings, if you want to change urls, edit `.env` file the URL bit.
otherwise settings can be passed directly in docker run as an environment variable.

```bash
$ docker run -e "URLS=https://www.google.com,https://www.facebook.com" -v ~/mrcrowley:/crawlie crawler
```

## How it works
the cron job runs once a week, every 7 days starting from the first run when the script is run.
it will crawl the urls and save them to a file in the folder you specified.
a folder with the current date will be created and the files will be saved there.
each URL will have its own file.