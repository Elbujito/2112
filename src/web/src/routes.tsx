import {
  MdDashboard,
  MdHome,
} from 'react-icons/md';

const routes = [
  {
    name: 'My 2112',
    layout: '/admin',
    path: 'default',
    icon: <MdHome className="text-inherit h-5 w-5" />,
    collapse: true,
    items: [
      {
        name: 'Satellite Tracking',
        layout: '/admin',
        path: '/default',
      },
      {
        name: 'World Map',
        layout: '/admin',
        path: '/world-map',
      },
    ],
  },
  {
    name: 'Game Management',
    path: '/admin',
    icon: <MdDashboard className="text-inherit h-5 w-5" />,
    collapse: true,
    items: [
      {
        name: 'New Game',
        layout: '/admin/games',
        path: '/new-game',
        exact: false,
      },
      {
        name: 'Overview',
        layout: '/admin/games',
        path: '/overview',
        exact: false,
      },
      {
        name: 'Reports',
        layout: '/admin/games',
        path: '/reports',
        exact: false,
      },
    ],
  },

];
export default routes;
