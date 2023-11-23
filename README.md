# Earthship

Collecting smart home data and serving it through a web frontend.  Also providing some forecasts based on external
influences and prior values.

## Data Sources

This is a work in progress, but these are the data sources that come to mind first.

1. SMA API for solar panel output
2. Xcel SDK for utilities usage (power, gas)
3. Nest API for house temperature, AC/heat usage
4. Rachio API for sprinkler usage
5. Some smart water meter (Flume makes one and has an API) to monitor overall water usage.
6. Cistern level monitoring on the rainwater tanks.
7. Some weather API (maybe just piggyback forecast data from Nest or Rachio?) for temp, sunlight, precipitation.
8. Others in the future?  Would love to categorize water usage.

The MVP will collect SMA, Nest, and Rachio usage.

## System Design

* Cron jobs to pull data from APIs and write to a relational database.  These jobs run for each API independently in
case they have different rate limits or data availability.  Each job hits the API for all available data and writes it
to one or more tables (this database will be suitably normalized).  The exact design of the database will depend on what
data are available, so I don't have a great idea yet.
* Cron job to train forecasts of every deliverable metric with enough data.
* Manual backfill job to fill missing data (useful for cold start or API outages).
* Backend to serve deliverables and predictions through an API.
* Frontend to display a web page with a dashboard of all deliverables.  Clicking on one opens a detailed page for just
that deliverable.

### Languages

The cron jobs can be written in any language.  I would like to experiment with something I don't already know. The model
training will be written in python.  The backend can be written in any language with a good API framework.  again I
would like to experiment with something new.  The frontend will be TS/react because I am not a web dev and that is the
only frontend language I know.

### Cloud design

Regardless of how this program is hosted, I am going to have a raspberry pi with a display panel in the house
permanently loaded to the frontend UI. I will do this in AWS initially and pay attention to the cost.  If the cost is
too high, then I will run this entirely from the pi.  The initial AWS architecture will run the cron jobs as lambdas
triggered by cron cloudwatch events.  The database will be an RDS instance, and the backend will be run with API
gateway.  The frontend will be hosted the same way as my personal website, which is hosted using route 53 and amplify.
