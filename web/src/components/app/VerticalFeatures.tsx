import dynamic from 'next/dynamic';
import { Section } from '../../shared/layout/Section';
import Skeleton from '@mui/material/Skeleton';

const VerticalFeaturesLoading = dynamic(() => import('../../shared/features/VerticalFeatureRow'), {
  loading: () => <Skeleton variant="rectangular" width={500} height={500}/>
})

export default function VerticalFeatures() {
  return (
    <Section
      styles={`py-6 max-w-screen-lg mx-auto px-3 `}
      title="Our services"
      description="2112 is a company that offers AI-powered gamification solutions to help businesses train and equip their employees to handle crises effectively. By using gamification, 2112 aims to make the learning process engaging, interactive for employees, making it easier for them to retain critical information and skills."
    >
      <VerticalFeaturesLoading 
        title="AI-powered solution"
        description="The AI-powered aspect of 2112's gamification solutions helps to personalize the training experience for each employee, tailoring the content and challenges to their unique needs and abilities. This approach can help ensure that each employee receives the training and support they need to handle crises effectively, whether they are dealing with a sudden security threat or a natural disaster."
        image="/assets/images/feature.png"
        imageAlt=""
      />
      <VerticalFeaturesLoading
        title="Real-time feedback"
        description="2112's gamification solutions can also track employee progress and provide real-time feedback, helping managers and supervisors to monitor employee readiness and identify areas that need improvement. This can be particularly helpful in identifying potential gaps in the workforce's crisis management skills and developing targeted training programs to address these issues."
        image="/assets/images/feature2.png"
        imageAlt=""
        reverse
      />
    </Section>
  );
}
