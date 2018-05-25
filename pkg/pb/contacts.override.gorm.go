package pb

import (
	"fmt"

	"github.com/infobloxopen/atlas-app-toolkit/rpc/errdetails"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BeforeToORM will add the primary e-mail to the list of e-mails if it isn't
// present already
func (m *Contact) BeforeToORM(ctx context.Context, c *ContactORM) error {
	if m.PrimaryEmail != "" {
		for _, mail := range m.Emails {
			if mail.Address == m.PrimaryEmail {
				return nil
			}
		}
		c.Emails = append(c.Emails, &EmailORM{Address: m.PrimaryEmail, IsPrimary: true})
	}
	return nil
}

// AfterToPB copies the primary e-mail address from the DB to the special PB field
func (m *ContactORM) AfterToPB(ctx context.Context, c *Contact) error {
	if len(m.Emails) == 0 {
		return nil
	}
	// find the primary e-mail in list of e-mails from DB
	for _, addr := range m.Emails {
		if addr != nil && addr.IsPrimary {
			c.PrimaryEmail = addr.Address
			break
		}
	}
	return nil
}

// Hooks:
// For more info: http://gorm.io/docs/hooks.html

// AfterFind is an example of a hook that runs after the Read function is complete
func (m *ContactORM) AfterFind() (err error) {
	// Business logic based on application
	// fmt.Println("Contact result:", m)
	return
}

// CustomList allows a custom error to be returned as a error response.
func (m *ContactsDefaultServer) CustomList(ctx context.Context, req *ListContactsRequest) (*ListContactsResponse, error) {
	res, err := DefaultListContact(ctx, m.DB)
	fmt.Println(res)
	if err != nil {
		st := status.Newf(codes.Internal, "Unable to Get contacts list.  Error %v", err)
		st, _ = st.WithDetails(errdetails.New(codes.Internal, "There was problem getting the contacts list", "CustomList"))
		return nil, st.Err()
	}
	return &ListContactsResponse{Results: res}, nil
}

// CustomErrorExample shows how to add a custom error with multiple details.
func (m *ContactsDefaultServer) CustomErrorExample(ctx context.Context, req *ErrorContactRequest) (*ErrorContactResponse, error) {
	st := status.Newf(codes.Internal, "Example error message")
	st, _ = st.WithDetails(errdetails.New(codes.Internal, "Example error message 1", "Example target 1"))
	st, _ = st.WithDetails(errdetails.New(codes.Internal, "Example error message 2", "Example target 2"))
	return nil, st.Err()
}
