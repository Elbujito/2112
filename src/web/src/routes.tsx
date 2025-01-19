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
      {
        name: 'Contexts',
        path: '/admin',
        icon: <MdDashboard className="text-inherit h-5 w-5" />,
        collapse: true,
        items: [
          {
            name: 'New Context',
            layout: '/admin/contexts',
            path: '/new-context',
            exact: false,
          },
          {
            name: 'Overview',
            layout: '/admin/contexts',
            path: '/overview',
            exact: false,
          },
          {
            name: 'Reports',
            layout: '/admin/contexts',
            path: '/reports',
            exact: false,
          },
        ],
      },
    ],
  },

];
export default routes;
