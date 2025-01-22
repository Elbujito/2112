import { LandingPrimaryVideoCtaSection } from 'components/landing/cta/LandingPrimaryCta';
import { LandingSocialProof } from 'components/landing/social-proof/LandingSocialProof';
import videoSrc from '/public/videos/landing.mp4';

import { Button } from 'components/shared/ui/button';
import { colors } from 'components/data/config/colors';

export default function Component() {
  const avatarItems = [
    {
      imageSrc: 'https://picsum.photos/id/64/100/100',
      name: 'John Doe',
    },
    {
      imageSrc: 'https://picsum.photos/id/65/100/100',
      name: 'Jane Doe',
    },
    {
      imageSrc: 'https://picsum.photos/id/669/100/100',
      name: 'Alice Doe',
    },
  ];

  return (
    <>
      <LandingPrimaryVideoCtaSection
        title="Master the Art of Space Warfare"
        description="2112 provides cutting-edge tools for space warfare assessment and strategy. Analyze scenarios, simulate outcomes, and develop superior tactics with our AI-driven platform, designed for the challenges of modern warfare in space."
        textPosition="center"
        videoPosition="center"
        videoSrc={videoSrc}
        withBackground
        variant="secondary"
        leadingComponent={
          <p className="font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-blue-500 via-blue-400 to-blue-500">
            The Future of Warfare in the Final Frontier
          </p>
        }
      >

        <div className="w-full mt-6 flex justify-center gap-4">
          <Button size="xl" className="p-7 text-xl" variant="outlineSecondary" asChild>
            <a href="/home/default">Start free today</a>
          </Button>
        </div>

        <LandingSocialProof
          className="w-full mt-6 justify-center"
          showRating
          numberOfUsers={25000}
          suffixText="satellite enthusiasts"
          avatarItems={avatarItems}
          size="large"
          disableAnimation
        />
      </LandingPrimaryVideoCtaSection>

      <div
        className="fixed top-0 left-0 w-full h-full -z-10"
        style={{
          backgroundImage: `url('data:image/svg+xml;utf8,${encodeURIComponent(
            ` <svg xmlns="http://www.w3.org/2000/svg"><defs><radialGradient id="a" cx="50%" cy="56.6%" r="50%" fx="50%" fy="56.6%" gradientUnits="userSpaceOnUse"><stop offset="0%" style="stop-color:${colors.primary.dark};stop-opacity:0.1"/><stop offset="54.99%" style="stop-color:${colors.primary.darker};stop-opacity:0.1"/><stop offset="100%" style="stop-color:${colors.secondary.darker};stop-opacity:0.1"/></radialGradient></defs><rect width="100%" height="100%" fill="url(#a)"/></svg>`,
          )}')`,
          backgroundSize: 'cover',
          backgroundRepeat: 'no-repeat',
          backgroundAttachment: 'fixed',
        }}
      ></div>
    </>
  );
}
