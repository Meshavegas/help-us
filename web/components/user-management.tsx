'use client';

import { useState, useEffect } from 'react';
import { apiClient, useApi, User, Address, Payment, Resource } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { 
  Sheet, 
  SheetContent, 
  SheetDescription, 
  SheetHeader, 
  SheetTitle, 
  SheetTrigger 
} from '@/components/ui/sheet';
import { 
  Users, 
  Eye, 
  Edit, 
  Trash2, 
  MoreHorizontal, 
  Search,
  UserPlus,
  Mail,
  Phone,
  MapPin,
  CreditCard,
  BookOpen
} from 'lucide-react';

interface UserManagementProps {
  currentUser: User;
}

export default function UserManagement({ currentUser }: UserManagementProps) {
  const [users, setUsers] = useState<User[]>([]);
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [userAddresses, setUserAddresses] = useState<Address[]>([]);
  const [userPayments, setUserPayments] = useState<Payment[]>([]);
  const [userResources, setUserResources] = useState<Resource[]>([]);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState<number | null>(null);

  const { loading, error, execute } = useApi<User[]>();
  const { loading: detailsLoading, execute: executeDetails } = useApi();
  const { loading: updateLoading, execute: executeUpdate } = useApi<User>();
  const { loading: deleteLoading, execute: executeDelete } = useApi();

  // Charger la liste des utilisateurs
  useEffect(() => {
    loadUsers();
  }, []);

  // Filtrer les utilisateurs selon le terme de recherche
  useEffect(() => {
    if (searchTerm) {
      const filtered = users.filter(user =>
        user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (user.famille?.family_name && user.famille.family_name.toLowerCase().includes(searchTerm.toLowerCase())) ||
        (user.enseignant?.specialization && user.enseignant.specialization.toLowerCase().includes(searchTerm.toLowerCase()))
      );
      setFilteredUsers(filtered);
    } else {
      setFilteredUsers(users);
    }
  }, [searchTerm, users]);

  const loadUsers = async () => {
    await execute(() => apiClient.getAllUsers());
  };

  const loadUserDetails = async (userId: number) => {
    // Charger les adresses
    await executeDetails(() => apiClient.getUserAddresses(userId));
    // Charger les paiements
    await executeDetails(() => apiClient.getUserPayments(userId));
    // Charger les ressources
    await executeDetails(() => apiClient.getUserResources(userId));
  };

  const handleViewUser = async (user: User) => {
    setSelectedUser(user);
    setIsDetailsOpen(true);
    await loadUserDetails(user.id);
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setIsEditOpen(true);
  };

  const handleUpdateUser = async (userData: any) => {
    if (!editingUser) return;

    await executeUpdate(() => apiClient.updateUser(editingUser.id, userData));
    setIsEditOpen(false);
    setEditingUser(null);
    await loadUsers();
  };

  const handleDeleteUser = async (userId: number) => {
    await executeDelete(() => apiClient.deleteUser(userId));
    setDeleteConfirm(null);
    await loadUsers();
  };

  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case 'administrator':
        return 'bg-red-100 text-red-800';
      case 'enseignant':
        return 'bg-blue-100 text-blue-800';
      case 'famille':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getRoleLabel = (role: string) => {
    switch (role) {
      case 'administrator':
        return 'Administrateur';
      case 'enseignant':
        return 'Enseignant';
      case 'famille':
        return 'Famille';
      default:
        return role;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Gestion des Utilisateurs</h2>
          <p className="text-muted-foreground">
            Gérez tous les utilisateurs de la plateforme éducative
          </p>
        </div>
        <Button>
          <UserPlus className="mr-2 h-4 w-4" />
          Nouvel Utilisateur
        </Button>
      </div>

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Recherche et Filtres</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center space-x-2">
            <Search className="h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Rechercher par nom, email, famille ou spécialisation..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="max-w-sm"
            />
          </div>
        </CardContent>
      </Card>

      {/* Users Table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Users className="mr-2 h-5 w-5" />
            Utilisateurs ({filteredUsers.length})
          </CardTitle>
          <CardDescription>
            Liste de tous les utilisateurs enregistrés
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : error ? (
            <div className="text-center py-8 text-red-600">
              Erreur: {error}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Utilisateur</TableHead>
                  <TableHead>Rôle</TableHead>
                  <TableHead>Contact</TableHead>
                  <TableHead>Statut</TableHead>
                  <TableHead>Informations Spécifiques</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredUsers.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell>
                      <div>
                        <div className="font-medium">{user.username}</div>
                        <div className="text-sm text-muted-foreground">ID: {user.id}</div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge className={getRoleBadgeColor(user.role)}>
                        {getRoleLabel(user.role)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="space-y-1">
                        <div className="flex items-center text-sm">
                          <Mail className="mr-1 h-3 w-3" />
                          {user.email}
                        </div>
                        {user.phone_number && (
                          <div className="flex items-center text-sm text-muted-foreground">
                            <Phone className="mr-1 h-3 w-3" />
                            {user.phone_number}
                          </div>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={user.is_active ? "default" : "secondary"}>
                        {user.is_active ? 'Actif' : 'Inactif'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {user.famille && (
                        <div className="text-sm">
                          <span className="font-medium">Famille:</span> {user.famille.family_name}
                        </div>
                      )}
                      {user.enseignant && (
                        <div className="text-sm">
                          <span className="font-medium">Spécialisation:</span> {user.enseignant.specialization}
                        </div>
                      )}
                      {user.administrator && (
                        <div className="text-sm text-muted-foreground">
                          Compte administrateur
                        </div>
                      )}
                    </TableCell>
                    <TableCell className="text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" className="h-8 w-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleViewUser(user)}>
                            <Eye className="mr-2 h-4 w-4" />
                            Voir détails
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleEditUser(user)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Modifier
                          </DropdownMenuItem>
                          {user.id !== currentUser.id && (
                            <DropdownMenuItem 
                              onClick={() => setDeleteConfirm(user.id)}
                              className="text-red-600"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Supprimer
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* User Details Sheet */}
      <Sheet open={isDetailsOpen} onOpenChange={setIsDetailsOpen}>
        <SheetContent className="w-[600px] sm:w-[800px]">
          <SheetHeader>
            <SheetTitle>Détails de l'utilisateur</SheetTitle>
            <SheetDescription>
              Informations complètes sur l'utilisateur sélectionné
            </SheetDescription>
          </SheetHeader>
          
          {selectedUser && (
            <div className="mt-6">
              <Tabs defaultValue="general" className="w-full">
                <TabsList className="grid w-full grid-cols-4">
                  <TabsTrigger value="general">Général</TabsTrigger>
                  <TabsTrigger value="addresses">Adresses</TabsTrigger>
                  <TabsTrigger value="payments">Paiements</TabsTrigger>
                  <TabsTrigger value="resources">Ressources</TabsTrigger>
                </TabsList>
                
                <TabsContent value="general" className="space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle>Informations Générales</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <label className="text-sm font-medium">Nom d'utilisateur</label>
                          <p className="text-sm text-muted-foreground">{selectedUser.username}</p>
                        </div>
                        <div>
                          <label className="text-sm font-medium">Email</label>
                          <p className="text-sm text-muted-foreground">{selectedUser.email}</p>
                        </div>
                        <div>
                          <label className="text-sm font-medium">Téléphone</label>
                          <p className="text-sm text-muted-foreground">
                            {selectedUser.phone_number || 'Non renseigné'}
                          </p>
                        </div>
                        <div>
                          <label className="text-sm font-medium">Rôle</label>
                          <p className="text-sm text-muted-foreground">
                            {getRoleLabel(selectedUser.role)}
                          </p>
                        </div>
                      </div>
                      
                      {selectedUser.famille && (
                        <div>
                          <label className="text-sm font-medium">Nom de famille</label>
                          <p className="text-sm text-muted-foreground">
                            {selectedUser.famille.family_name}
                          </p>
                        </div>
                      )}
                      
                      {selectedUser.enseignant && (
                        <div className="space-y-2">
                          <div>
                            <label className="text-sm font-medium">Spécialisation</label>
                            <p className="text-sm text-muted-foreground">
                              {selectedUser.enseignant.specialization}
                            </p>
                          </div>
                          <div>
                            <label className="text-sm font-medium">Qualifications</label>
                            <p className="text-sm text-muted-foreground">
                              {selectedUser.enseignant.qualifications}
                            </p>
                          </div>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                </TabsContent>
                
                <TabsContent value="addresses" className="space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center">
                        <MapPin className="mr-2 h-4 w-4" />
                        Adresses
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      {detailsLoading ? (
                        <div className="text-center py-4">Chargement...</div>
                      ) : userAddresses.length > 0 ? (
                        <div className="space-y-3">
                          {userAddresses.map((address) => (
                            <div key={address.id} className="border rounded-lg p-3">
                              <div className="flex items-center justify-between">
                                <div>
                                  <p className="font-medium">{address.street}</p>
                                  <p className="text-sm text-muted-foreground">
                                    {address.city}, {address.postal_code}, {address.country}
                                  </p>
                                </div>
                                {address.is_primary && (
                                  <Badge>Principale</Badge>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-center text-muted-foreground py-4">
                          Aucune adresse enregistrée
                        </p>
                      )}
                    </CardContent>
                  </Card>
                </TabsContent>
                
                <TabsContent value="payments" className="space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center">
                        <CreditCard className="mr-2 h-4 w-4" />
                        Paiements
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      {detailsLoading ? (
                        <div className="text-center py-4">Chargement...</div>
                      ) : userPayments.length > 0 ? (
                        <div className="space-y-3">
                          {userPayments.map((payment) => (
                            <div key={payment.id} className="border rounded-lg p-3">
                              <div className="flex items-center justify-between">
                                <div>
                                  <p className="font-medium">
                                    {payment.amount} {payment.currency}
                                  </p>
                                  <p className="text-sm text-muted-foreground">
                                    {payment.description || 'Paiement'}
                                  </p>
                                  <p className="text-xs text-muted-foreground">
                                    {new Date(payment.payment_date).toLocaleDateString()}
                                  </p>
                                </div>
                                <Badge>{payment.status}</Badge>
                              </div>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-center text-muted-foreground py-4">
                          Aucun paiement enregistré
                        </p>
                      )}
                    </CardContent>
                  </Card>
                </TabsContent>
                
                <TabsContent value="resources" className="space-y-4">
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center">
                        <BookOpen className="mr-2 h-4 w-4" />
                        Ressources
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      {detailsLoading ? (
                        <div className="text-center py-4">Chargement...</div>
                      ) : userResources.length > 0 ? (
                        <div className="space-y-3">
                          {userResources.map((resource) => (
                            <div key={resource.id} className="border rounded-lg p-3">
                              <div className="flex items-center justify-between">
                                <div>
                                  <p className="font-medium">{resource.title}</p>
                                  <p className="text-sm text-muted-foreground">
                                    {resource.description}
                                  </p>
                                  <p className="text-xs text-muted-foreground">
                                    Type: {resource.type}
                                  </p>
                                </div>
                                <Badge variant={resource.is_active ? "default" : "secondary"}>
                                  {resource.is_active ? 'Actif' : 'Inactif'}
                                </Badge>
                              </div>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-center text-muted-foreground py-4">
                          Aucune ressource accessible
                        </p>
                      )}
                    </CardContent>
                  </Card>
                </TabsContent>
              </Tabs>
            </div>
          )}
        </SheetContent>
      </Sheet>

      {/* Delete Confirmation Dialog */}
      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-semibold mb-4">Confirmer la suppression</h3>
            <p className="text-muted-foreground mb-6">
              Êtes-vous sûr de vouloir supprimer cet utilisateur ? Cette action est irréversible.
            </p>
            <div className="flex justify-end space-x-2">
              <Button 
                variant="outline" 
                onClick={() => setDeleteConfirm(null)}
                disabled={deleteLoading}
              >
                Annuler
              </Button>
              <Button 
                variant="destructive" 
                onClick={() => handleDeleteUser(deleteConfirm)}
                disabled={deleteLoading}
              >
                {deleteLoading ? 'Suppression...' : 'Supprimer'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
} 