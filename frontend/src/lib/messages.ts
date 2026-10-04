export const errorMessages = {
  pins: {
    fetch: "Failed to fetch pins",
    create: "Failed to create pin",
    update: "Failed to update pin",
    delete: "Failed to delete pin",
  },
  notes: {
    create: "Failed to create note",
    update: "Failed to update note",
    delete: "Failed to delete note",
  },
  plan: {
    create: "Failed to create plan",
    fetch: "Failed to fetch plan",
  },
  spots: {
    search: "Failed to search spots",
  },
  user: {
    register: "Failed to register user",
    update: "Failed to update user",
    delete: "Failed to delete user",
    loginStart: "Failed to start login",
    loginVerify: "Failed to verify login",
    fetchCurrentUser: "Failed to fetch current user",
    logout: "Failed to log out",
  },
  countryGuides: {
    list: "Failed to fetch country guides",
    get: "Failed to fetch country guide",
  },
  planMembers: {
    list: "Failed to fetch plan members",
    add: "Failed to add plan member",
    remove: "Failed to remove plan member",
  },
} as const;
