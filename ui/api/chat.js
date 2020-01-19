import { getUser as User } from './user';
// chat menu
const Menu = [
  {
    text: 'Chat',
    icon: 'chat',
    to: { path: '/chat/messaging' },
  }
];

const getChatById = (uuid) => {
  return (uuid !== undefined) ? Groups.find(x => x.uuid === uuid) : Groups[0];
};

export {
  Menu,
  getChatById
};
