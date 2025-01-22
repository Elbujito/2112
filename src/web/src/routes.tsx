import {
  MdDashboard,
  MdHome,
} from 'react-icons/md';
import ProfileIcon from 'components/icons/ProfileIcon';
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
    name: 'Games Management',
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
  {
    name: 'Users Management',
    path: '/admin',
    icon: <ProfileIcon />,
    collapse: true,
    items: [
      {
        name: 'Overview',
        layout: '/admin/users',
        path: '/overview',
        exact: false,
      },
    ],
  },

];
export default routes;
