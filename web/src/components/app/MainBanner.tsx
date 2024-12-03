import { CTABanner } from '../../shared/cta/CTABanner';
import { Section } from '../../shared/layout/Section';

export default function MainBanner() {
  return (
    <div>
      <Section styles={`py-2 max-w-screen-lg mx-auto px-3`}>
        <CTABanner />
      </Section>
    </div>
  );
}
