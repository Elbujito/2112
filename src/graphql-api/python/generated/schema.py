"""Generated GraphQL Schema"""

type Query {
  satellitePosition(id: ID!): SatellitePosition
  satelliteTle(id: ID!): SatelliteTle
  satellitePositionsInRange(id: ID!, startTime: String!, endTime: String!): [SatellitePosition!]!
}

type SatellitePosition {
  id: ID!
  name: String!
  latitude: Float!
  longitude: Float!
  altitude: Float!
  timestamp: String!
}

type SatelliteTle {
  id: ID!
  name: String!
  tleLine1: String!
  tleLine2: String!
}

type Subscription {
  satellitePositionUpdated(id: ID!): SatellitePosition
}