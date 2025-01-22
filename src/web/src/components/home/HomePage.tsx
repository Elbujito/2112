'use client';

import Logo from './specta/Logo';

import OneVideoCta from 'components/home/specta/bricks/1-video-cta';
import ThreeTestimonialInline from 'components/home/specta/bricks/3-testimonial-inline';
import FourProductFeature from 'components/home/specta/bricks/4-product-feature';
import SixProductFeature from 'components/home/specta/bricks/6-product-feature';
import NineSaleCta from 'components/home/specta/bricks/9-sale-cta';
import LandingHeader from './specta/HomeHeader';
import Footer from 'components/footer/Footer';

export const Specta = () => {
    return (
        <div className="flex flex-col w-full min-h-screen">
            <header className="relative z-50">
                <LandingHeader
                    className="absolute top-0 w-full"
                    logo={<Logo className="h-9 w-auto" />}
                    logoDark={<Logo className="h-9 w-auto" />}
                />
            </header>

            <main className="flex flex-col items-center w-full space-y-16">
                <section className="w-full">
                    <OneVideoCta />
                </section>

                <section className="w-full">
                    <ThreeTestimonialInline />
                </section>

                <section className="w-full py-16">
                    <FourProductFeature />
                </section>

                <section className="w-full">
                    <SixProductFeature />
                </section>

                <section className="w-full py-16">
                    <NineSaleCta />
                </section>
            </main>

            <footer className="w-full mt-16">
                <Footer />
            </footer>
        </div>
    );
};
