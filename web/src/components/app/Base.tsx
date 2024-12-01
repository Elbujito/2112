import { Meta } from '../../shared/layout/Meta';
import { AppConfig } from '../../config/AppConfig';
import { Footer } from '../footer/Footer';
import { Hero } from '../hero/Hero';
import { Analytics } from '@vercel/analytics/react';
import VerticalFeatures from './VerticalFeatures';
import Tracker from '../tracker/tracker';


export default function Base() {

  return (
  <div className="antialiased text-white-100">
    <Meta title={AppConfig.title} description={AppConfig.description} />
    <Hero title="AI Gamification for Crisis-Ready Workforces"/>
    <Tracker/>
    <VerticalFeatures/>
    <Footer />
    <Analytics />
  </div>
  );

}

export { Base };
