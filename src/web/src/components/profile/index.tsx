'use client';

import React from 'react';
import { UserProfile } from '@clerk/clerk-react';

export const UserProfilePage = () => {
    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
            <div className="w-full max-w-4xl p-6 bg-white rounded-lg shadow-lg dark:bg-gray-800">
                <UserProfile
                    appearance={{
                        elements: {
                            card: 'shadow-lg border rounded-lg p-6 bg-white dark:bg-gray-800',
                            navbar: 'bg-blue-600 text-white dark:bg-blue-800',
                            button: 'rounded-md px-4 py-2 bg-blue-500 hover:bg-blue-400 text-white font-semibold',
                        },
                    }}
                    routing="path"
                    path="/user-profile"
                    additionalOAuthScopes={{
                        google: ['calendar', 'drive'],
                        github: ['repo'],
                    }}
                />
            </div>
        </div>
    );
};

export default UserProfilePage;
