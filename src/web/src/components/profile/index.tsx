'use client';

import React from 'react';
import { UserProfile, SignedIn, SignedOut, SignIn } from '@clerk/clerk-react';

export const UserProfilePage = () => {
    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-navy-900">
            <SignedIn>
                <div className="w-full max-w-4xl p-6 bg-white rounded-lg shadow-lg dark:bg-navy-900">
                    <UserProfile
                        additionalOAuthScopes={{
                            google: ['calendar', 'drive'],
                            github: ['repo'],
                        }}
                    />
                </div>
            </SignedIn>
            <SignedOut>
                <div className="text-center text-gray-800 dark:text-white">
                    <SignIn />
                </div>
            </SignedOut>
        </div>
    );
};

export default UserProfilePage;
