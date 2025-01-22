import { Button } from 'components/shared/ui/button';

import { LandingProductFeature } from 'components/landing/LandingProductFeature';
import { LandingProductFeatureKeyPoints } from 'components/landing/LandingProductFeatureKeyPoints';
import satelliteTrackingImg from '/public/img/satellites/satellites-operations.webp';

export default function Component() {
  const keyPoints = [
    {
      title: 'Real-Time Satellite Visibility',
      description:
        'Easily check which satellites are visible above any location, at any time, with a single click.',
    },
    {
      title: 'Global Coverage',
      description:
        'Access satellite data for any point on Earth, whether you’re planning a mission or stargazing.',
    },
    {
      title: 'Interactive Visualization',
      description:
        'Explore satellite trajectories in an intuitive, visually rich interface for better decision-making.',
    },
  ];

  return (
    <LandingProductFeature
      titleComponent={
        <>
          <p className="text-xl font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-blue-500 via-blue-400 to-blue-500">
            Explore
          </p>

          <h2 className="text-4xl font-semibold leading-tight">
            Discover visible satellites in real-time, from anywhere in the world
          </h2>
        </>
      }
      descriptionComponent={
        <>
          <p>
            Dive into the fascinating world of satellite tracking with our
            cutting-edge tools. Whether you're a professional or just curious,
            our platform provides all the data and insights you need to explore
            satellites in the sky, live and in real-time.
          </p>

          <LandingProductFeatureKeyPoints
            variant="secondary"
            keyPoints={keyPoints}
            className="mt-4"
          />

          <Button className="mt-8" variant="outlineSecondary" asChild>
            <a href="#">Start exploring for free</a>
          </Button>

          <p className="text-sm">No subscription required for basic features.</p>
        </>
      }
      // imageSrc={satelliteTrackingImg.src}
      imageAlt="Real-time satellite tracking interface"
      imagePosition="right"
      imagePerspective="none"
      zoomOnHover={false}
      withBackground
      variant="secondary"
    />
  );
}
