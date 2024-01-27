import {
    ArrayField,
    ArrayInput,
    Datagrid,
    List,
    NumberField,
    SearchInput, Show, ShowButton,
    SimpleFormIterator,
    SimpleShowLayout,
    TextField
} from 'react-admin';

import { CreateButton, ExportButton, TopToolbar } from 'react-admin';

const PostListActions = () => (
    <TopToolbar>
        <ShowButton />
    </TopToolbar>
);
export const RoleList = () => (

    <List
        filters={orderFilters}
        bulkActionButtons={false}
        sort={{ field: 'permissionCount', order: 'DESC' }}
        actions={<PostListActions />}
    >

        <Datagrid rowClick="show">
            <TextField label="Role" source="name" />
            <TextField source="title" />
            <TextField source="description" />
            <TextField source="stage" />
            <NumberField source="permissionCount" />
        </Datagrid>
    </List>
);

const orderFilters = [
    <SearchInput source="q" alwaysOn />,
];

import { useRecordContext } from 'react-admin';

const IncludedPermissionsField = () => {
    const record = useRecordContext();
    return (
        <ul>
            {record.includedPermissions.map((item: string) => (
                <li key={item}>{item}</li>
            ))}
        </ul>
    )
};

export const RoleShow = () => (
    <Show>
        <SimpleShowLayout>
            <TextField source="name" />
            <TextField source="title" />
            <TextField source="description" />
            <NumberField source="permissionCount" />
            <IncludedPermissionsField/>
        </SimpleShowLayout>
    </Show>
);