# wordle-clone-server-ARCHIVE

This here is an archive for the REST(ish) API I created for a uni project.

This ran on AWS Elastic Beanstalk and connected to a separate MySQL database running on an EC2 instance.

If you try running this locally, it will NOT work, unless you have a MySQL database with an identical schema and add it in the `./db_client/db_client.go` file! Then, you just run it as a normal Docker container.
