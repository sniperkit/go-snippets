
export function selectCommunities(state) {
  const { users } = state.entities;
  return Object.keys(users).filter(id => users[id].type == "community").map(id => users[id]);
}

// search through the state tree for
export function selectUserByUsername(state, username) {
  const { users } = state.entities;
  const userId = Object.keys(users).find(id => (users[id].username == username));
  return userId ? users[userId] : undefined;
}

export function selectUserById(state, id) {
  return state.entities.users[id];
}

export function selectCommunityUsers(state, id) {
  const { users } = state.entities;
  const { communityUsers } = state.pagination;

  if (communityUsers[id]) {
    return communityUsers[id].ids.map(id => users[id]);
  }

  return [];
}

export function selectLocalUser(state, id) {
  // const { users } = state.locals;
  // if (!users) {
  //   return undefined;
  // }
  // const u = Object.keys(users).map(k => users[k]).find(u)
  // return (u.length == 1) ? u[0] : undefined;
  return state.locals.users[id];
}
