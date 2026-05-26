import { createFileRoute } from '@tanstack/react-router'
import ProfileList from '../../components/profiles/ProfileList'

export const Route = createFileRoute('/profiles')({
    component: ProfilesComponent,
})

function ProfilesComponent() {
    return (
        <ProfileList />
    )
}