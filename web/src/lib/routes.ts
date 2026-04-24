export const ROUTES = {
  OTHERS: {
    NOT_FOUND: "/not-found",
  },
  MAIN: {
    AUTH: {
      SIGN_IN: "/sign-in",
      SIGN_UP: "/sign-up",
      COMPLETE_SETUP: "/complete-setup",
    },
    ASSETS: {
      RESOURCE: "/assets/",
      LIST: "/assets",
      SHOWCASE: "/assets/:id",
      EDIT: "/assets/:id/edit",
      CREATE: "/assets/create",
    },
    STUDIO: {
      RESOURCE: "/studio/",
      LIST: "/studio",
      SHOWCASE: "/studio/:id",
      EDIT: "/studio/:id/edit",
    }
  },
  ADMIN: {
    MAIN: "/admin",
  }
};
