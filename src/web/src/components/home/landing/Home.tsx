import React from 'react';
import Earth from './Earth';
import ThreeTestimonialInline from 'components/home/specta/bricks/3-testimonial-inline';
import FourProductFeature from 'components/home/specta/bricks/4-product-feature';
import SixProductFeature from 'components/home/specta/bricks/6-product-feature';
import NineSaleCta from 'components/home/specta/bricks/9-sale-cta';
import LandingHeader from '../specta/HomeHeader';
import Footer from 'components/footer/Footer';
import Logo from '../specta/Logo';
import { LandingPrimaryVideoCtaSection } from 'components/landing/cta/LandingPrimaryCta';
import { Button } from 'components/shared/ui/button';
import { colors } from 'components/data/config/colors';

const Home = () => {

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
        <div className="flex flex-col w-full min-h-screen">
            <header className="relative z-50">
                <LandingHeader
                    className="absolute top-0 w-full"
                    logo={<Logo className="h-9 w-auto" />}
                    logoDark={<Logo className="h-9 w-auto" />}
                />
            </header>

            <main className="flex flex-col items-center w-full space-y-16">
                <section className="w-full relative">
                    <div className="absolute inset-0 flex flex-col items-center justify-center">
                        <h1 className="text-white text-4xl font-bold"></h1>
                        <LandingPrimaryVideoCtaSection
                            title="Master the Art of Space Warfare"
                            description=""
                            textPosition="center"
                            videoPosition="center"
                            videoSrc={""}
                            withBackground
                            variant="secondary"
                            leadingComponent={
                                <p className="z-10 font-cursive font-semibold tracking-wider bg-clip-text bg-gradient-to-r text-transparent from-blue-500 via-blue-400 to-blue-500">
                                    The Future of Warfare in the Final Frontier
                                </p>
                            }
                        > <Button size="xl" className="p-7 mt-6 text-xl z-10" variant="outlineSecondary" asChild>
                                <a href="/home/default">Start free today</a>
                            </Button></LandingPrimaryVideoCtaSection>



                    </div>
                    <Earth />
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

            <footer className="w-full mt-16">
                <Footer />
            </footer>
        </div>
    );
};

export default Home;
