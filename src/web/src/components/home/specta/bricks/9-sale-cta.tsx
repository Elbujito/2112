import { LandingSaleCtaSection } from 'components/landing/cta/LandingSaleCta';
import { LandingSocialProof } from 'components/landing/social-proof/LandingSocialProof';

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

export default function Component() {
  return (
    <LandingSaleCtaSection
      titleComponent={
        <>
          <p className="text-xl font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-blue-500 via-blue-400 to-blue-500">
            It takes 1 minute
          </p>

          <h2 className="text-4xl font-semibold leading-tight">
            The faster, easier way to explore satellites
          </h2>
        </>
      }
      descriptionComponent={
        <>
          <p>
            Jump in today and discover how easy it is to track, explore, and
            learn about satellites orbiting Earth in real-time with Project 2112.
          </p>

          <LandingSocialProof
            className="w-full mt-6"
            showRating
            numberOfUsers={25000}
            suffixText="enthusiastic explorers"
            avatarItems={avatarItems}
            disableAnimation
          >
            <p className="text-xs">trusted by 25,000+ satellite enthusiasts</p>
          </LandingSocialProof>
        </>
      }
      ctaHref="#"
      ctaLabel="Start exploring now"
      withBackgroundGlow
      withBackground
    />
  );
}
