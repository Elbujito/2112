import React, { useEffect, useMemo } from "react";
import { useQuery, useSubscription, useMutation, gql } from "@apollo/client";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import VisibilityTimeline from "./VisibilityTimeline";

const GET_CACHED_VISIBILITIES = gql`
  query GetCachedSatelliteVisibilities(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    cachedSatelliteVisibilities(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    ) {
      satelliteId
      satelliteName
      aos
      los
      userLocation {
        latitude
        longitude
        radius
        horizon
      }
    }
  }
`;

const VISIBILITY_SUBSCRIPTION = gql`
  subscription OnVisibilityUpdated(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    satelliteVisibilityUpdated(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    ) {
      satelliteId
      satelliteName
      aos
      los
      userLocation {
        latitude
        longitude
        radius
        horizon
      }
    }
  }
`;

const REQUEST_VISIBILITIES_MUTATION = gql`
  mutation RequestSatelliteVisibilities(
    $uid: String!
    $userLocation: UserLocationInput!
    $startTime: String!
    $endTime: String!
  ) {
    requestSatelliteVisibilities(
      uid: $uid
      userLocation: $userLocation
      startTime: $startTime
      endTime: $endTime
    )
  }
`;

interface VisibilityTimelineWithDataProps {
    uid: string; // User ID
    userLocation: { latitude: number; longitude: number }; // User's location
}

const VisibilityTimelineWithData: React.FC<VisibilityTimelineWithDataProps> = ({ uid, userLocation }) => {
    const radius = 500; // 500 km radius
    const horizon = 30; // 30 degrees horizon angle

    // Set up time range
    const startTime = useMemo(() => new Date().toISOString(), []);
    const endTime = useMemo(() => new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(), []);

    const userLocationInput = {
        ...userLocation,
        radius,
        horizon,
        uid,
    };

    // Mutation to request satellite visibilities
    const [requestVisibilities, { loading: mutationLoading, error: mutationError }] = useMutation(
        REQUEST_VISIBILITIES_MUTATION
    );

    // Query to fetch initial cached data
    const { data, loading: queryLoading, error: queryError } = useQuery(GET_CACHED_VISIBILITIES, {
        variables: { uid, userLocation: userLocationInput, startTime, endTime },
        skip: !uid,
    });

    // Subscription for real-time updates
    const { data: subscriptionData, error: subscriptionError } = useSubscription(VISIBILITY_SUBSCRIPTION, {
        variables: { uid, userLocation: userLocationInput, startTime, endTime },
        skip: !uid,
    });

    // Call the mutation to request visibilities on component mount
    useEffect(() => {
        if (uid) {
            requestVisibilities({
                variables: { uid, userLocation: userLocationInput, startTime, endTime },
            }).catch((err) => {
                console.error("Error requesting satellite visibilities:", err);
            });
        }
    }, [uid, userLocationInput, startTime, endTime, requestVisibilities]);

    // Combine query and subscription data
    const visibilities = useMemo(() => {
        const queriedData = data?.cachedSatelliteVisibilities || [];
        const liveUpdates = subscriptionData?.satelliteVisibilityUpdated || [];
        return [...queriedData, ...liveUpdates];
    }, [data, subscriptionData]);

    // Handle loading state
    if (queryLoading || mutationLoading) {
        return (
            <Center>
                <Spinner />
                <Text ml={2}>Loading visibilities...</Text>
            </Center>
        );
    }

    // Handle error states
    if (queryError || subscriptionError || mutationError) {
        console.error("Query Error:", queryError);
        console.error("Subscription Error:", subscriptionError);
        console.error("Mutation Error:", mutationError);
        return (
            <Box>
                <Text color="red.500">
                    Error loading visibilities:
                    {queryError?.message || subscriptionError?.message || mutationError?.message}
                </Text>
            </Box>
        );
    }

    // Map data for the timeline
    const timelineData = visibilities.map((visibility) => ({
        current: false,
        day: new Date(visibility.aos).getDate().toString(),
        weekday: new Date(visibility.aos).toLocaleDateString(undefined, { weekday: "short" }),
        hours: `${new Date(visibility.aos).toLocaleTimeString()} - ${new Date(visibility.los).toLocaleTimeString()}`,
        title: "Visibility Window",
        satellite: visibility.satelliteName || "Unknown Satellite",
        noradID: visibility.satelliteId || "N/A",
    }));

    return <VisibilityTimeline location={userLocation} data={timelineData} />;
};

export default VisibilityTimelineWithData;
