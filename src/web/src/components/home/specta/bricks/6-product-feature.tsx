import { Button } from 'components/shared/ui/button';
import { LandingProductFeature } from 'components/landing/LandingProductFeature';
import { LandingProductFeatureKeyPoints } from 'components/landing/LandingProductFeatureKeyPoints';
import satelliteThreatImg from '/public/img/satellites/satellites-threats.webp';

export default function Component() {
  const keyPoints = [
    {
      title: 'Threat Detection',
      description:
        'Monitor real-time alerts for potential satellite collisions or hazardous events.',
    },
    {
      title: 'Global Activity',
      description: 'Analyze satellite operations and risks across countries.',
    },
    {
      title: 'Threat Leaderboard',
      description:
        'Track countries generating the most alerts and identify patterns in space activities.',
    },
  ];

  return (
    <LandingProductFeature
      titleComponent={
        <>
          <p className="text-xl font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-red-500 via-red-400 to-red-500">
            Stay Alert
          </p>

          <h2 className="text-4xl font-semibold leading-tight">
            Real-time Satellite Threat Alerts and Leaderboard
          </h2>
        </>
      }
      descriptionComponent={
        <>
          <p>
            Gain critical insights into satellite activities and potential
            threats in space. Our platform provides live alerts for satellite
            collisions, operational risks, and other hazardous events. Stay
            ahead of emerging threats with a detailed country leaderboard and
            global trend analysis.
          </p>

          <LandingProductFeatureKeyPoints
            variant="secondary"
            keyPoints={keyPoints}
            className="mt-4"
          />

          <Button className="mt-8" variant="outlineSecondary" asChild>
            <a href="#">Discover our offers</a>
          </Button>

          <p className="text-sm mt-2">
            Access essential features at no cost. Premium plans available for
            advanced tracking.
          </p>

          <div className="mt-6">
            <h3 className="text-lg font-semibold">Threat Leaderboard</h3>
            <p className="text-sm">
              Keep track of countries generating the most satellite threat
              alerts. Understand global space activity trends and assess risks
              with real-time data.
            </p>
          </div>
        </>
      }
      // imageSrc={satelliteThreatImg.src}
      imageAlt="Satellite threat detection interface"
      imagePosition="right"
      imagePerspective="none"
      zoomOnHover={false}
      withBackground
      variant="secondary"
    />
  );
}
