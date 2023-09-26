'use strict'

const autocannon = require('autocannon')

const validPayload = {
  subject: "Foo Bar",
  created: "2023-09-26T09:00:00.000-04:00",
  body: "A test post.",
  author: {
    name: "John Doe",
    email: "jdoe@example.com"
  }
}
const validBodyConfig = {
  url: 'http://127.0.0.1:8080/message/new',
  connections: 30,
  duration: 30,
  workers: 3,
  method: 'POST',
  headers: {
    'content-type': 'application/json'
  },
  body: JSON.stringify(validPayload)
}

const invalidPayload = {
  subject: "Foo Bar",
  created: "2023-09-26",
  body: "A test post.",
  author: {
    name: "John Doe",
    email: "jdoe@example.com ~~ is invalid"
  }
}
const invalidBodyConfig = {
  url: 'http://127.0.0.1:8080/message/new',
  connections: 30,
  duration: 30,
  workers: 3,
  method: 'POST',
  headers: {
    'content-type': 'application/json'
  },
  body: JSON.stringify(invalidPayload)
}

main()
  .catch(error => {
    console.error(error)
  })
  .then(() => {})

async function main() {
  let instance
  process.once('SIGINT', () => {
    instance.stop()
  })

  console.log("running valid payload test...")
  instance = autocannon(validBodyConfig)
  autocannon.track(instance, {
    renderResultsTable: false,
    renderLatencyTable: false
  })
  let results = await instance
  printResults(results)

  console.log("\n\nrunning invalid payload test...")
  instance = autocannon(invalidBodyConfig)
  autocannon.track(instance, {
    renderResultsTable: false,
    renderLatencyTable: false
  })
  results = await instance
  printResults(results)
}

function printResults(results) {
  // console.dir(results)
  console.table({
    statuses: {
      '1xx': results['1xx'],
      '2xx': results['2xx'],
      '3xx': results['3xx'],
      '4xx': results['4xx'],
      '5xx': results['5xx']
    }
  })
  console.table({
    requests: {
      average: results.requests.average,
      mean: results.requests.mean,
      stddev: results.requests.stddev,
      total: results.requests.total
    }
  })
}
