import { LandingTestimonialInline } from 'components/landing/testimonial/LandingTestimonialInline';
import { LandingTestimonialInlineItem } from 'components/landing/testimonial/LandingTestimonialInlineItem';

export default function Component() {
  return (
    <LandingTestimonialInline>
      <LandingTestimonialInlineItem
        name="John Doe"
        text="Project 2112 has completely transformed how we track and analyze satellite data."
        suffix="Aerospace Engineer at NASA"
      />

      <LandingTestimonialInlineItem
        name="Jane Doe"
        text="The best platform for real-time satellite tracking and insights."
        suffix="Astronomy Enthusiast"
      />

      <LandingTestimonialInlineItem
        name="Alice Doe"
        text="I've discovered hundreds of satellites and learned so much about orbital mechanics."
        suffix="Professor of Astrophysics"
      />

      <LandingTestimonialInlineItem
        name="Guido Ross"
        text="Automating satellite data workflows has never been easier. Project 2112 is a game-changer!"
        suffix="Satellite Systems Analyst at SpaceX"
      />
    </LandingTestimonialInline>
  );
}
