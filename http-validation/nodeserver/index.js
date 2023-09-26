'use strict'

const server = require('fastify')({
  ajv: {
    plugins: [require('ajv-formats')]
  }
})

const payloadSchema = {
  type: 'object',
  properties: {
    subject: { type: 'string' },
    created: {
      type: 'string',
      format: 'date-time'
    },
    body: { type: 'string' },
    author: {
      type: 'object',
      properties: {
        name: { type: 'string' },
        email: {
          type: 'string',
          format: 'email'
        }
      },
      required: ['name', 'email']
    }
  },
  required: ['subject', 'created', 'body', 'author']
}

server.route({
  method: 'post',
  path: '/message/new',
  schema: {
    body: payloadSchema
  },
  async handler() {
    return "ok"
  }
})

server.listen({
  port: '8080',
  host: process.env['HOST_ADDR'] || '127.0.0.1'
})
