// DATA='{"any": "data"}' ENDPOINT=api/endpoint METHOD=verb ACCESS_TOKEN=yourAccessToken,omitempty k6 run --vus 100 --iterations 100000 stress-test/generic.js

import http from 'k6/http';

import { check } from 'k6';

const host = __ENV.HTTP_HOST || 'http://localhost';
const post = __ENV.HTTP_PORT || '8080';
const authorization = __ENV.ACCESS_TOKEN ? `Bearer ${__ENV.ACCESS_TOKEN}` : '';
const data = __ENV.DATA || undefined
const method = __ENV.METHOD || 'get'

export default function () {
  let res = http[method](
    `${host}:${post}/${__ENV.ENDPOINT}`,
    data,
    {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': authorization,
      },
    },
  );

  check(res, { 'status was 200': (r) => r.status == 200 });
}
