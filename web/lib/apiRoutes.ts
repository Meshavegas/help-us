export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

export type Role =
  | 'public'        // No authentication required
  | 'authenticated' // Any authenticated user (admin, teacher, parent, child)
  | 'admin'
  | 'teacher'
  | 'parent'
  | 'child'

export interface ApiEndpoint {
  method: HttpMethod
  path: string
  roles: Role[]
  description?: string
}

/**
 * Complete list of Go-backend API endpoints and their minimal roles.
 */
export const apiEndpoints: ApiEndpoint[] = [
  { method: 'GET',  path: '/health',                          roles: ['public'],          description: 'Health-check' },
  { method: 'POST', path: '/api/v1/auth/register',            roles: ['public'],          description: 'Register' },
  { method: 'POST', path: '/api/v1/auth/login',               roles: ['public'],          description: 'Login' },
  { method: 'POST', path: '/api/v1/auth/refresh',             roles: ['public'],          description: 'Refresh token' },
  { method: 'POST', path: '/api/v1/auth/logout',              roles: ['authenticated'],   description: 'Logout' },

  { method: 'GET',  path: '/api/v1/profile',                  roles: ['authenticated'],   description: 'Current user profile' },
  { method: 'PUT',  path: '/api/v1/profile',                  roles: ['authenticated'],   description: 'Update profile' },

  { method: 'GET',  path: '/api/v1/users/:id',                roles: ['authenticated'],   description: 'User details' },
  { method: 'GET',  path: '/api/v1/users/:id/addresses',      roles: ['authenticated'],   description: 'User addresses' },
  { method: 'GET',  path: '/api/v1/users/:id/payments',       roles: ['authenticated'],   description: 'User payments' },
  { method: 'GET',  path: '/api/v1/users/:id/resources',      roles: ['authenticated'],   description: 'User resources' },
  { method: 'GET',  path: '/api/v1/users',                    roles: ['admin'],           description: 'List users' },
  { method: 'PUT',  path: '/api/v1/users/:id',                roles: ['admin'],           description: 'Update user' },
  { method: 'DELETE',path: '/api/v1/users/:id',               roles: ['admin'],           description: 'Delete user' },
  { method: 'GET',  path: '/api/v1/admin/users',              roles: ['admin'],           description: 'Admin list users' },
  { method: 'PUT',  path: '/api/v1/admin/users/:id',          roles: ['admin'],           description: 'Admin update user' },
  { method: 'DELETE',path: '/api/v1/admin/users/:id',         roles: ['admin'],           description: 'Admin delete user' },

  { method: 'GET',  path: '/api/v1/teacher/courses',          roles: ['teacher','admin'], description: 'Teacher courses' },
  { method: 'GET',  path: '/api/v1/family/missions',          roles: ['parent','admin'],  description: 'Family missions' },

  { method: 'GET',  path: '/api/v1/familles',                 roles: ['admin'],           description: 'List familles' },
  { method: 'GET',  path: '/api/v1/familles/:id',             roles: ['authenticated'],   description: 'Famille details' },
  { method: 'PUT',  path: '/api/v1/familles/:id',             roles: ['authenticated'],   description: 'Update famille' },
  { method: 'DELETE',path: '/api/v1/familles/:id',            roles: ['admin'],           description: 'Delete famille' },
  { method: 'GET',  path: '/api/v1/familles/:id/teachers',    roles: ['authenticated'],   description: 'Famille teachers' },
  { method: 'GET',  path: '/api/v1/familles/:id/missions',    roles: ['authenticated'],   description: 'Famille missions' },
  { method: 'GET',  path: '/api/v1/familles/:id/courses',     roles: ['authenticated'],   description: 'Famille courses' },
  { method: 'GET',  path: '/api/v1/familles/:id/payments',    roles: ['authenticated'],   description: 'Famille payments' },
  { method: 'POST', path: '/api/v1/familles/:id/reviews',     roles: ['authenticated'],   description: 'Add famille review' },
  { method: 'GET',  path: '/api/v1/familles/:id/options',     roles: ['authenticated'],   description: 'Famille options' },

  { method: 'GET',  path: '/api/v1/missions',                 roles: ['authenticated'],   description: 'List missions' },
  { method: 'POST', path: '/api/v1/missions',                 roles: ['authenticated'],   description: 'Create mission' },
  { method: 'GET',  path: '/api/v1/missions/:id',             roles: ['authenticated'],   description: 'Mission details' },
  { method: 'PUT',  path: '/api/v1/missions/:id',             roles: ['authenticated'],   description: 'Update mission' },
  { method: 'DELETE',path: '/api/v1/missions/:id',            roles: ['authenticated'],   description: 'Delete mission' },
  { method: 'GET',  path: '/api/v1/missions/:id/courses',     roles: ['authenticated'],   description: 'Mission courses' },
  { method: 'GET',  path: '/api/v1/missions/:id/reports',     roles: ['authenticated'],   description: 'Mission reports' },
  { method: 'GET',  path: '/api/v1/missions/:id/payments',    roles: ['authenticated'],   description: 'Mission payments' },
  { method: 'PUT',  path: '/api/v1/missions/:id/stop',        roles: ['authenticated'],   description: 'Stop mission' },
  { method: 'PUT',  path: '/api/v1/missions/:id/extend',      roles: ['authenticated'],   description: 'Extend mission' },

  { method: 'GET',  path: '/api/v1/courses',                  roles: ['authenticated'],   description: 'List courses' },
  { method: 'POST', path: '/api/v1/courses',                  roles: ['authenticated'],   description: 'Create course' },
  { method: 'GET',  path: '/api/v1/courses/:id',              roles: ['authenticated'],   description: 'Course details' },
  { method: 'PUT',  path: '/api/v1/courses/:id',              roles: ['authenticated'],   description: 'Update course' },
  { method: 'DELETE',path: '/api/v1/courses/:id',             roles: ['authenticated'],   description: 'Delete course' },
  { method: 'PUT',  path: '/api/v1/courses/:id/schedule',     roles: ['authenticated'],   description: 'Schedule course' },
  { method: 'PUT',  path: '/api/v1/courses/:id/cancel',       roles: ['authenticated'],   description: 'Cancel course' },
  { method: 'PUT',  path: '/api/v1/courses/:id/complete',     roles: ['authenticated'],   description: 'Complete course' },
  { method: 'POST', path: '/api/v1/courses/:id/declare',      roles: ['authenticated'],   description: 'Declare course' },
  { method: 'GET',  path: '/api/v1/courses/:id/payments',     roles: ['authenticated'],   description: 'Course payments' },

  { method: 'GET',  path: '/api/v1/enseignants',              roles: ['authenticated'],   description: 'List teachers' },
  { method: 'POST', path: '/api/v1/enseignants',              roles: ['admin'],           description: 'Create teacher' },
  { method: 'GET',  path: '/api/v1/enseignants/:id',          roles: ['authenticated'],   description: 'Teacher details' },
  { method: 'PUT',  path: '/api/v1/enseignants/:id',          roles: ['authenticated'],   description: 'Update teacher' },
  { method: 'DELETE',path: '/api/v1/enseignants/:id',         roles: ['admin'],           description: 'Delete teacher' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/students', roles: ['authenticated'],   description: 'Teacher students' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/missions', roles: ['authenticated'],   description: 'Teacher missions' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/courses',  roles: ['authenticated'],   description: 'Teacher courses' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/payments', roles: ['authenticated'],   description: 'Teacher payments' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/reports',  roles: ['authenticated'],   description: 'Teacher reports' },
  { method: 'GET',  path: '/api/v1/enseignants/:id/options',  roles: ['authenticated'],   description: 'Teacher options' },
  { method: 'GET',  path: '/api/v1/enseignants/nearby',       roles: ['authenticated'],   description: 'Nearby teachers' },

  { method: 'GET',  path: '/api/v1/offers',                   roles: ['authenticated'],   description: 'List offers' },
  { method: 'POST', path: '/api/v1/offers',                   roles: ['authenticated'],   description: 'Create offer' },
  { method: 'GET',  path: '/api/v1/offers/:id',               roles: ['authenticated'],   description: 'Offer details' },
  { method: 'PUT',  path: '/api/v1/offers/:id',               roles: ['authenticated'],   description: 'Update offer' },
  { method: 'DELETE',path: '/api/v1/offers/:id',              roles: ['authenticated'],   description: 'Delete offer' },
  { method: 'GET',  path: '/api/v1/offers/:id/options',       roles: ['authenticated'],   description: 'Offer options' },
  { method: 'PUT',  path: '/api/v1/offers/:id/close',         roles: ['authenticated'],   description: 'Close offer' },
  { method: 'GET',  path: '/api/v1/offers/active',            roles: ['authenticated'],   description: 'Active offers' },
  { method: 'GET',  path: '/api/v1/offers/search',            roles: ['authenticated'],   description: 'Search offers' },

  { method: 'GET',  path: '/api/v1/options',                  roles: ['authenticated'],   description: 'List options' },
  { method: 'POST', path: '/api/v1/options',                  roles: ['authenticated'],   description: 'Create option' },
  { method: 'GET',  path: '/api/v1/options/:id',              roles: ['authenticated'],   description: 'Option details' },
  { method: 'PUT',  path: '/api/v1/options/:id',              roles: ['authenticated'],   description: 'Update option' },
  { method: 'DELETE',path: '/api/v1/options/:id',             roles: ['authenticated'],   description: 'Delete option' },
  { method: 'PUT',  path: '/api/v1/options/:id/accept',       roles: ['authenticated'],   description: 'Accept option' },
  { method: 'PUT',  path: '/api/v1/options/:id/decline',      roles: ['authenticated'],   description: 'Decline option' },
  { method: 'PUT',  path: '/api/v1/options/:id/cancel',       roles: ['authenticated'],   description: 'Cancel option' },
  { method: 'GET',  path: '/api/v1/options/pending',          roles: ['authenticated'],   description: 'Pending options' },
  { method: 'GET',  path: '/api/v1/options/expiring',         roles: ['authenticated'],   description: 'Expiring options' },

  { method: 'GET',  path: '/api/v1/addresses',                roles: ['admin'],           description: 'List addresses' },
  { method: 'POST', path: '/api/v1/addresses',                roles: ['authenticated'],   description: 'Create address' },
  { method: 'GET',  path: '/api/v1/addresses/:id',            roles: ['authenticated'],   description: 'Address details' },
  { method: 'PUT',  path: '/api/v1/addresses/:id',            roles: ['authenticated'],   description: 'Update address' },
  { method: 'DELETE',path: '/api/v1/addresses/:id',           roles: ['authenticated'],   description: 'Delete address' },
  { method: 'GET',  path: '/api/v1/addresses/geocode',        roles: ['authenticated'],   description: 'Geocode address' },
  { method: 'GET',  path: '/api/v1/addresses/route',          roles: ['authenticated'],   description: 'Calculate route' },
]

export function getEndpointsForRole(role: Role | 'public'): ApiEndpoint[] {
  if (role === 'public') {
    return apiEndpoints.filter(ep => ep.roles.includes('public'))
  }

  return apiEndpoints.filter(ep =>
    ep.roles.includes(role) ||
    ep.roles.includes('authenticated') ||
    (role === 'admin' && ep.roles.includes('admin'))
  )
} 