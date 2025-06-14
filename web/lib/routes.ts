import { 
  LayoutDashboard, 
  Users, 
  GraduationCap, 
  BookOpen, 
  Settings,
  Calendar,
  MessageSquare,
  FileText,
  CreditCard,
  Building2,
  MapPin
} from "lucide-react"

export type Route = {
  title: string
  href: string
  icon: any
  roles: string[]
  children?: Route[]
}

export const routes: Route[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
    roles: ["admin", "teacher", "student", "family"]
  },
  {
    title: "Users",
    href: "/users",
    icon: Users,
    roles: ["admin"],
    children: [
      {
        title: "All Users",
        href: "/users",
        icon: Users,
        roles: ["admin"]
      },
      {
        title: "Teachers",
        href: "/users/teachers",
        icon: GraduationCap,
        roles: ["admin"]
      },
      {
        title: "Students",
        href: "/users/students",
        icon: Users,
        roles: ["admin"]
      },
      {
        title: "Families",
        href: "/users/families",
        icon: Building2,
        roles: ["admin"]
      }
    ]
  },
  {
    title: "Courses",
    href: "/courses",
    icon: BookOpen,
    roles: ["admin", "teacher", "student", "family"],
    children: [
      {
        title: "All Courses",
        href: "/courses",
        icon: BookOpen,
        roles: ["admin", "teacher", "student", "family"]
      },
      {
        title: "Schedule",
        href: "/courses/schedule",
        icon: Calendar,
        roles: ["admin", "teacher", "student", "family"]
      },
      {
        title: "My Courses",
        href: "/courses/my-courses",
        icon: BookOpen,
        roles: ["teacher", "student"]
      }
    ]
  },
  {
    title: "Messages",
    href: "/messages",
    icon: MessageSquare,
    roles: ["admin", "teacher", "student", "family"]
  },
  {
    title: "Documents",
    href: "/documents",
    icon: FileText,
    roles: ["admin", "teacher", "student", "family"]
  },
  {
    title: "Payments",
    href: "/payments",
    icon: CreditCard,
    roles: ["admin", "family"],
    children: [
      {
        title: "All Payments",
        href: "/payments",
        icon: CreditCard,
        roles: ["admin"]
      },
      {
        title: "My Payments",
        href: "/payments/my-payments",
        icon: CreditCard,
        roles: ["family"]
      }
    ]
  },
  {
    title: "Locations",
    href: "/locations",
    icon: MapPin,
    roles: ["admin"],
    children: [
      {
        title: "All Locations",
        href: "/locations",
        icon: MapPin,
        roles: ["admin"]
      },
      {
        title: "Add Location",
        href: "/locations/add",
        icon: MapPin,
        roles: ["admin"]
      }
    ]
  },
  {
    title: "Settings",
    href: "/settings",
    icon: Settings,
    roles: ["admin", "teacher", "student", "family"]
  }
] 