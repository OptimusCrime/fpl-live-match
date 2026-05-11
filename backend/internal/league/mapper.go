package league

func toPlayerOwnersDTO(po *PlayerOwners) playerOwnersDTO {
	owners := make([]playerOwnerDTO, len(po.Owners))
	for i, o := range po.Owners {
		owners[i] = playerOwnerDTO{EntryName: o.EntryName, PlayerName: o.PlayerName, Rank: o.Rank}
	}
	return playerOwnersDTO{PlayerID: po.PlayerID, Owners: owners, TotalEntries: po.TotalEntries}
}
