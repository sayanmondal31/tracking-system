import http from "k6/http";
import { Trend, Counter } from "k6/metrics";
import { sleep } from "k6";

// Define custom metrics for our benchmark
const rateLimitTime = new Trend("rate_limit_check_time_us"); // Time in microseconds
const status200 = new Counter("status_200");
const status429 = new Counter("status_429");

// Test Configuration
export const options = {
  scenarios: {
    constant_load: {
      executor: "constant-arrival-rate",
      rate: 1000000, // target rate: 5000 requests per second
      timeUnit: "1s",
      duration: "30s", // run the test for 30 seconds
      preAllocatedVUs: 500, // allocate enough virtual users up front
      maxVUs: 1000,
    },
  },
};

export default function () {
  const res = http.get("http://localhost:3000/health");

  // 1. Log Status Codes
  if (res.status === 200) {
    status200.add(1);
  } else if (res.status === 429) {
    status429.add(1);
  }

  // 2. Parse and log rate limit check time from headers
  const headerTime = res.headers["X-Ratelimit-Duration-Us"];
  if (headerTime) {
    rateLimitTime.add(parseInt(headerTime, 10));
  }
}
